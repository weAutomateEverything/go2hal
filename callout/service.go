package callout

import (

	"net/http"
	"io/ioutil"
	"strings"

	"fmt"

	"errors"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/snmp"
	"github.com/zamedic/go2hal/jira"
	"github.com/zamedic/go2hal/config"
)


type Service interface {
	/*
	InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
	*/
	InvokeCallout(title, message string)

	getFirstCallName() (string, error)

}

type service struct {
	alert alert.Service
	snmp snmp.Service
	jira jira.Service
}

func NewService(alert alert.Service,snmp snmp.Service,jira jira.Service) Service{
	return &service{alert,snmp,jira}
}

/*
InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
 */
func (s *service)InvokeCallout(title, message string) {

	s.alert.SendError(fmt.Errorf("invoking callout for: %s, %s", title, message))
	s.snmp.SendSNMPMessage()
	n, err := s.getFirstCallName()
	if err != nil {
		s.alert.SendError(err)
		n = "DEFAULT"
	}
	s.jira.CreateJira(title, message,n)
}



func (s *service)getFirstCallName() (string, error) {
	c, err := s.configStore.GetCalloutDetails()
	if err != nil {
		return "", err
	}
	if c == nil {
		return "", errors.New("no callout set")
	}
	endpoint := c.URL
	if endpoint == "" {
		return "DEFAULT", nil
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		s.alert.SendError(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.alert.SendError(err)
		return "", err
	}
	bodyString := string(body)
	split := strings.SplitAfter(bodyString, "<font color='red' size=2>")
	names := strings.Split(split[1], "</font>")
	if len(names) == 0 {
		return "", errors.New("no callout found")
	}
	return names[0], nil

}


