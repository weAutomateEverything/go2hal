package httpSmoke

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NewService(alertService alert.Service, store Store, callout callout.Service, checkCount, ErrorCount metrics.Counter) Service {
	s := &service{
		alert:      alertService,
		store:      store,
		callout:    callout,
		checkCount: checkCount,
		errorCount: ErrorCount,
		notifySSLExpiryTimeout: time.Now(),
	}
	go func() {
		s.monitorEndpoints()
	}()
	return s
}

type Service interface {
	setTimeOut(minutes int64)
	getEndpoints(group uint32) ([]httpEndpoint, error)
	addHttpEndpoint(ctx context.Context, name, url, method string, parameters []parameters, threshold int, chat uint32) error
	deleteEndpoint(id string, chat uint32)
	sendSSLExpiryAlert(alertFrequency int, expiryStatus string, ctx context.Context, endpoint httpEndpoint)
}

type service struct {
	alert      alert.Service
	store      Store
	callout    callout.Service
	timeout    time.Time
	timeoutSet bool

	checkCount metrics.Counter
	errorCount metrics.Counter

	notifySSLExpiryTimeout time.Time
}

func (s *service) getEndpoints(group uint32) ([]httpEndpoint, error) {
	return s.store.getHTMLEndpointsByChat(group)
}

func (s *service) addHttpEndpoint(ctx context.Context, name, url, method string, parameters []parameters, threshold int, chat uint32) (err error) {
	v := httpEndpoint{
		Chat:       chat,
		Name:       name,
		Endpoint:   url,
		Method:     method,
		Parameters: parameters,
		Threshold:  threshold,
	}

	err = s.checkEndpoint(ctx, v)
	if err != nil {
		return
	}

	s.alert.SendAlert(ctx, chat, emoji.Sprintf(":new: Successfully added endpoint %v %v. The bot will now alert you once the checks fails %v times in succession. ", name, url, threshold))

	return s.store.addHTMLEndpoint(v)
}

func (s *service) deleteEndpoint(id string, chat uint32) {
	panic("implement me")
}

func (s *service) checkEndpoint(ctx context.Context, endpoint httpEndpoint) error {
	response, err := s.doHTTPEndpoint(ctx, endpoint)
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

	ctx := context.Background()

	response, err := s.doHTTPEndpoint(ctx, endpoint)
	s.checkCount.With("endpoint", endpoint.Name).Add(1)
	if err != nil {
		msg := emoji.Sprintf(":smoking: :x: *Smoke Test  Alert*\nName: %s \nEndpoint: %s \nError: %s", endpoint.Name,
			endpoint.Endpoint, err.Error())
		s.checkAlert(ctx, endpoint, msg)
		s.errorCount.With("endpoint", endpoint.Name).Add(1)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		s.errorCount.With("endpoint", endpoint.Name).Add(1)
		msg, _ := ioutil.ReadAll(response.Body)
		error := emoji.Sprintf(":smoking: :x: *HTTP Alert*\nName: %s \nEndpoint: %s \nDid not receive a 200 success "+
			"response code. Received %d response code. Body Message %s", endpoint.Name, endpoint.Endpoint,
			response.StatusCode, msg)
		s.checkAlert(ctx, endpoint, error)
		return
	}
	if !endpoint.Passing && endpoint.ErrorCount >= endpoint.Threshold {
		s.alert.SendAlert(ctx, endpoint.Chat, emoji.Sprintf(":smoking: :white_check_mark: smoke test %s back to normal", endpoint.Name))

	}

	if err := s.store.successfulEndpointTest(&endpoint); err != nil {
		s.alert.SendError(ctx, err)
	}
	if response.TLS != nil && len(response.TLS.PeerCertificates) != 0 {
		certExpiry := response.TLS.PeerCertificates[0].NotAfter
		daysTillExpiry := s.daysToExpiry(certExpiry)
		expiryStatus, alertFrequency := s.confirmCertExpiry(certExpiry, endpoint.Endpoint, daysTillExpiry)

		if expiryStatus != "" {
			s.sendSSLExpiryAlert(alertFrequency, expiryStatus, ctx, endpoint)
		}
	}
}

func (s *service) checkAlert(ctx context.Context, endpoint httpEndpoint, msg string) {
	if err := s.store.failedEndpointTest(&endpoint, msg); err != nil {
		s.alert.SendError(ctx, err)
	}
	if endpoint.Threshold > 0 {
		if endpoint.Threshold == endpoint.ErrorCount {
			msg = fmt.Sprintf("We have dectected a problem with HTTP Endpoint %v, it has failed %v times in a row. The error is %v", endpoint.Name, endpoint.Threshold, msg)
			s.callout.InvokeCallout(ctx, endpoint.Chat, fmt.Sprintf("Some Test failur)es for %s", endpoint.Name), msg, true)

		}
		if endpoint.ErrorCount >= endpoint.Threshold {
			s.checkTimeout(ctx, endpoint.Chat, msg)
		}
	}
}

func (s *service) checkTimeout(ctx context.Context, chat uint32, msg string) {
	if !s.timeoutSet || time.Now().Local().After(s.timeout) {
		s.alert.SendAlert(ctx, chat, msg)

		if s.timeoutSet {
			s.alert.SendAlert(ctx, chat, emoji.Sprintf(":alarm_clock: - Smoke Alerts expired. The bot will now be sending alerts for smoke failures again"))

			s.timeoutSet = false
		}
	}
}

func (s service) doHTTPEndpoint(ctx context.Context, endpoint httpEndpoint) (resp *http.Response, err error) {

	var request *http.Request
	switch method := strings.ToUpper(endpoint.Method); method {
	case "POST":
		if len(endpoint.Parameters) > 1 {
			return nil, errors.New("POST expects 0 or 1 parameters to pass as a body, normally a json string. To send a form, please use method POST_FORM")
		}
		body := ""
		if len(endpoint.Parameters) == 1 {
			body = endpoint.Parameters[0].Value
		}

		request, err = http.NewRequest("POST", endpoint.Endpoint, bytes.NewBufferString(body))
		request.Header.Add("Content-Type", "application/json")
	case "POST_FORM":
		values := url.Values{}
		for _, value := range endpoint.Parameters {
			values.Add(value.Name, value.Value)
		}
		request, err = http.NewRequest("POST", endpoint.Endpoint, nil)
		request.Form = values

	default:
		request, err = http.NewRequest("GET", endpoint.Endpoint, nil)
	}
	if err != nil {
		return
	}
	c := &http.Client{Transport: defaultTransport}
	resp, err = c.Do(request)
	if err != nil {
		return
	}

	return

}

func (s *service) setTimeOut(minutes int64) {
	s.timeout = time.Now().Local().Add(time.Minute * time.Duration(minutes))
	s.timeoutSet = true
}

func (s *service) daysToExpiry(expiryDate time.Time) float64 {

	duration := expiryDate.Sub(time.Now())
	return math.Floor(duration.Hours() / 24)
}

func (s *service) confirmCertExpiry(expiryDate time.Time, endpoint string, expiryDays float64) (string, int) {

	if expiryDays <= 54 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 1
	}
	if expiryDays <= 60 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 1
	}
	if expiryDays <= 70 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 2
	}
	if expiryDays <= 85 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 2
	}
	if expiryDays <= 100 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 3
	}
	if expiryDays <= 120 {
		return emoji.Sprintf(":rotating_light: SSL certificate for %v expires withing %v days!\nExpiry date: %v", endpoint, expiryDays, expiryDate), 3
	}
	return "", 0
}

func (s *service) sendSSLExpiryAlert(alertFrequency int, expiryStatus string, ctx context.Context, endpoint httpEndpoint) {
	if alertFrequency == 3 && time.Now().After(s.notifySSLExpiryTimeout.Add(3 * time.Hour)){
		s.alert.SendAlert(ctx, endpoint.Chat, expiryStatus)
		s.notifySSLExpiryTimeout = time.Now()
	}else if alertFrequency == 2 && time.Now().After(s.notifySSLExpiryTimeout.Add(2 * time.Hour)){
		s.alert.SendAlert(ctx, endpoint.Chat, expiryStatus)
		s.notifySSLExpiryTimeout = time.Now()
	}else if alertFrequency == 1 && time.Now().After(s.notifySSLExpiryTimeout.Add(1 * time.Hour)){
		s.alert.SendAlert(ctx, endpoint.Chat, expiryStatus)
		s.notifySSLExpiryTimeout = time.Now()
	}
}

var defaultTransport http.RoundTripper = &http.Transport{
	Proxy: nil,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
