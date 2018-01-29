package http


import (
	"time"
	"net/http"
	"fmt"
	"log"
	"strings"
	"errors"
	"bytes"
	"net/url"
	"io/ioutil"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/callout"
)

type Service interface {

}

type service struct {
	alert alert.Service
	store Store
	callout callout.Service
}

func NewService(alertService alert.Service, store Store,callout callout.Service) Service {
	s := &service{alertService, store, callout}
	go func() {
		s.monitorEndpoints()
	}()
	return s
}


func (s service)checkEndpoint(endpoint httpEndpoint) error {
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

func (s service)monitorEndpoints() {
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

func (s service)checkHTTP(endpoint httpEndpoint) {
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
		s.alert.SendAlert(emoji.Sprintf(":smoking: :white_check_mark: smoke test %s back to normal", endpoint.Name))
	}

	if err := s.store.successfulEndpointTest(&endpoint); err != nil {
		s.alert.SendError(err)
	}
}

func (s service)checkAlert(endpoint httpEndpoint, msg string) {
	if err := s.store.failedEndpointTest(&endpoint, msg); err != nil {
		s.alert.SendError(err)
	}
	s.alert.SendError(errors.New(msg))
	if endpoint.Threshold > 0 {
		if endpoint.Threshold == endpoint.ErrorCount {
			s.callout.InvokeCallout(fmt.Sprintf("Some Test failures for %s", endpoint.Name),msg)
		}
		if endpoint.ErrorCount >= endpoint.Threshold {
			s.alert.SendAlert(msg)
		}
	}
}

func (s service)doHTTPEndpoint(endpoint httpEndpoint) (*http.Response, error) {
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

