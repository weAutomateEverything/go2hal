package callout

import (
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"

	"errors"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/jira"
	"github.com/zamedic/go2hal/snmp"
	"golang.org/x/net/context"
	"os"
)

type Service interface {
	/*
		InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
	*/
	InvokeCallout(ctx context.Context, title, message string)

	getFirstCallName(ctx context.Context) (string, error)
}

type service struct {
	alert alert.Service
	snmp  snmp.Service
	jira  jira.Service
}

func NewService(alert alert.Service, snmp snmp.Service, jira jira.Service) Service {
	return &service{alert, snmp, jira}
}

/*
InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
*/
func (s *service) InvokeCallout(ctx context.Context, title, message string) {

	s.alert.SendError(ctx, fmt.Errorf("invoking callout for: %s, %s", title, message))
	s.snmp.SendSNMPMessage(ctx)
	n, err := s.getFirstCallName(ctx)
	if err != nil {
		s.alert.SendError(ctx, err)
		n = "DEFAULT"
	}
	s.jira.CreateJira(ctx, title, message, n)
}

func (s *service) getFirstCallName(ctx context.Context) (string, error) {
	endpoint := getCalloutDetails()
	if endpoint == "" {
		return "DEFAULT", nil
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		s.alert.SendError(ctx, err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.alert.SendError(ctx, err)
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
func getCalloutDetails() string {
	return os.Getenv("CALLOUT_URL")
}
