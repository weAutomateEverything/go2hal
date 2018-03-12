package callout

import (
	"net/http"
	"strings"

	"fmt"

	"errors"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/jira"
	"github.com/weAutomateEverything/go2hal/snmp"
	"golang.org/x/net/context"
	"gopkg.in/xmlpath.v2"
	"os"
)

type Service interface {
	/*
		InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
	*/
	InvokeCallout(ctx context.Context, title, message string)

	getFirstCall(ctx context.Context) (name string, number string, err error)
}

type service struct {
	alert alert.Service
	snmp  snmp.Service
	jira  jira.Service
	alexa halaws.Service
}

func NewService(alert alert.Service, snmp snmp.Service, jira jira.Service, alexa halaws.Service) Service {
	return &service{alert, snmp, jira, alexa}
}

/*
InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
*/
func (s *service) InvokeCallout(ctx context.Context, title, message string) {

	s.alert.SendError(ctx, fmt.Errorf("invoking callout for: %s, %s", title, message))
	s.snmp.SendSNMPMessage(ctx)
	name, phone, err := s.getFirstCall(ctx)
	if err != nil {
		s.alert.SendError(ctx, err)
		name = "DEFAULT"
	}
	s.jira.CreateJira(ctx, title, message, name)
	if s.alexa != nil {
		s.alexa.SendAlert(phone)
	}
}

func (s *service) getFirstCall(ctx context.Context) (name string, number string, err error) {
	endpoint := getCalloutDetails()
	if endpoint == "" {
		return "DEFAULT", "", nil
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		s.alert.SendError(ctx, err)
		return "DEFAULT", "", err
	}
	nodes, err := xmlpath.ParseHTML(resp.Body)

	if err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error decoding html for callout list. %v", err))
		return "DEFAULT", "", err
	}

	namePath := xmlpath.MustCompile("/html/body/div[2]/fieldset[1]/table/tbody/tr[1]/th/font")
	phonePath := xmlpath.MustCompile("/html/body/div[2]/fieldset[1]/table/tbody/tr[3]/td[1]")

	name, ok := namePath.String(nodes)
	if !ok {
		s.alert.SendAlert(ctx, "unable to retrieve first call user from callout portal")
		return "DEFAULT", "", errors.New("unable to retrieve first call user from portal due ot xml parsing issue. ")
	}

	phone, ok := phonePath.String(nodes)
	if !ok {
		s.alert.SendAlert(ctx, "unable to retrieve phone number from callout portal")
		return "DEFAULT", "", errors.New("unable to retrieve phone number from portal due ot xml parsing issue. ")
	}

	phone = strings.Replace(phone, "-", "", -1)
	phone = strings.Replace(phone, " ", "", -1)
	phone = strings.Replace(phone, "0", "+27", 1)

	return name, phone, nil

}
func getCalloutDetails() string {
	return os.Getenv("CALLOUT_URL")
}
