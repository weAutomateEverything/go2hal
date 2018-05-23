package httpSmoke

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewService(alertService alert.Service, store Store, callout callout.Service) Service {
	s := &service{alert: alertService, store: store, callout: callout}
	go func() {
		s.monitorEndpoints()
	}()
	return s
}

type Service interface {
	setTimeOut(minutes int64)
	getEndpoints(group uint32) ([]httpEndpoint, error)
	addHttpEndpoint(name, url, method string, parameters []parameters, threshold int, chat uint32) error
	deleteEndpoint(id string, chat uint32)
}

type service struct {
	alert      alert.Service
	store      Store
	callout    callout.Service
	timeout    time.Time
	timeoutSet bool
}

func (s *service) getEndpoints(group uint32) ([]httpEndpoint, error) {
	return s.store.getHTMLEndpointsByChat(group)
}

func (s *service) addHttpEndpoint(name, url, method string, parameters []parameters, threshold int, chat uint32) (err error) {
	v := httpEndpoint{
		Chat:       chat,
		Name:       name,
		Endpoint:   url,
		Method:     method,
		Parameters: parameters,
		Threshold:  threshold,
	}

	err = s.checkEndpoint(v)
	if err != nil {
		return
	}

	return s.store.addHTMLEndpoint(v)
}

func (s *service) deleteEndpoint(id string, chat uint32) {
	panic("implement me")
}

func (s *service) checkEndpoint(endpoint httpEndpoint) error {
	response, err := s.doHTTPEndpoint(endpoint)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("response code recevied %d", response.StatusCode)
	}
	return nil
}

func (s *service) monitorEndpoints() {
	log.Println("Starting HTTP Endpoint monitor")
	for true {
		endpoints := s.store.getHTMLEndpoints()
		if endpoints != nil {
			for _, endpoint := range endpoints {
				s.checkHTTP(endpoint)
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func (s *service) checkHTTP(endpoint httpEndpoint) {
	response, err := s.doHTTPEndpoint(endpoint)
	if err != nil {
		msg := emoji.Sprintf(":smoking: :x: *Smoke Test  Alert*\nName: %s \nEndpoint: %s \nError: %s", endpoint.Name,
			endpoint.Endpoint, err.Error())
		s.checkAlert(endpoint, msg)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		msg, _ := ioutil.ReadAll(response.Body)
		error := emoji.Sprintf(":smoking: :x: *HTTP Alert*\nName: %s \nEndpoint: %s \nDid not receive a 200 success "+
			"response code. Received %d response code. Body Message %s", endpoint.Name, endpoint.Endpoint,
			response.StatusCode, msg)
		s.checkAlert(endpoint, error)
		return
	}
	if !endpoint.Passing && endpoint.ErrorCount >= endpoint.Threshold {
		s.alert.SendAlert(context.TODO(), endpoint.Chat, emoji.Sprintf(":smoking: :white_check_mark: smoke test %s back to normal", endpoint.Name))

	}

	if err := s.store.successfulEndpointTest(&endpoint); err != nil {
		s.alert.SendError(context.TODO(), err)
	}
	if response.TLS != nil && len(response.TLS.PeerCertificates) != 0 {
		certExpiry := response.TLS.PeerCertificates[0].NotAfter
		daysTillExpiry := s.daysToExpiry(certExpiry)
		expiryStatus := s.confirmCertExpiry(certExpiry, endpoint.Endpoint, daysTillExpiry)

		if expiryStatus != "" {
			err := errors.New(expiryStatus)
			s.alert.SendError(context.TODO(), err)
		}
	}
}

func (s *service) checkAlert(endpoint httpEndpoint, msg string) {
	if err := s.store.failedEndpointTest(&endpoint, msg); err != nil {
		s.alert.SendError(context.TODO(), err)
	}
	if endpoint.Threshold > 0 {
		if endpoint.Threshold == endpoint.ErrorCount {
			s.callout.InvokeCallout(context.TODO(), endpoint.Chat, fmt.Sprintf("Some Test failures for %s", endpoint.Name), msg, nil)

		}
		if endpoint.ErrorCount >= endpoint.Threshold {
			s.checkTimeout(endpoint.Chat, msg)
		}
	}
}

func (s *service) checkTimeout(chat uint32, msg string) {
	if !s.timeoutSet || time.Now().Local().After(s.timeout) {
		s.alert.SendAlert(context.TODO(), chat, msg)

		if s.timeoutSet {
			s.alert.SendAlert(context.TODO(), chat, emoji.Sprintf(":alarm_clock: - Smoke Alerts expired. The bot will now be sending alerts for smoke failures again"))

			s.timeoutSet = false
		}
	}
}

func (s service) doHTTPEndpoint(endpoint httpEndpoint) (*http.Response, error) {
	switch method := strings.ToUpper(endpoint.Method); method {
	case "POST":
		if len(endpoint.Parameters) > 1 {
			return nil, errors.New("POST expects 0 or 1 parameters to pass as a body, normally a json string. To send a form, please use method POST_FORM")
		}
		body := ""
		if len(endpoint.Parameters) == 1 {
			body = endpoint.Parameters[0].Value
		}
		return http.Post(endpoint.Endpoint, "application/json", bytes.NewBufferString(body))
	case "POST_FORM":
		values := url.Values{}
		for _, value := range endpoint.Parameters {
			values.Add(value.Name, value.Value)
		}
		return http.PostForm(endpoint.Endpoint, values)
	default:
		return http.Get(endpoint.Endpoint)
	}

}

func (s *service) setTimeOut(minutes int64) {
	s.timeout = time.Now().Local().Add(time.Minute * time.Duration(minutes))
	s.timeoutSet = true
}

func (s *service) daysToExpiry(expiryDate time.Time) float64 {

	duration := expiryDate.Sub(time.Now())
	return math.Floor(duration.Hours() / 24)
}

func (s *service) confirmCertExpiry(expiryDate time.Time, endpoint string, expiryDays float64) string {

	if expiryDays <= 54 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 54 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	if expiryDays <= 60 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 60 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	if expiryDays <= 70 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 70 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	if expiryDays <= 85 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 85 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	if expiryDays <= 100 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 100 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	if expiryDays <= 120 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing 120 days!\nExpiry date: %v", endpoint, expiryDate)
	}
	return ""
}
