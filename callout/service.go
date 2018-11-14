// Package callout provides a mechanism to invoke various different forms of callout depending on the services
// linked when creating the service.
package callout

import (
	"fmt"
	"time"

	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/jira"
	"github.com/weAutomateEverything/go2hal/snmp"
	"golang.org/x/net/context"
)

//Service interface for the Callout Service
type Service interface {
	//InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
	InvokeCallout(ctx context.Context, chat uint32, title, message string, ack bool) error
}

type service struct {
	alert     alert.Service
	snmp      snmp.Service
	jira      jira.Service
	alexa     halaws.Service
	firstcall firstCall.Service
	store     Store
}

// NewService creates a new Callout Service. Parameters can be passed in as Nil should they not be required.
// any items that are nil will simply not be invoked.
func NewService(alert alert.Service, firstcall firstCall.Service, snmp snmp.Service, jira jira.Service, alexa halaws.Service, store Store) Service {
	s := &service{
		snmp:      snmp,
		alert:     alert,
		alexa:     alexa,
		firstcall: firstcall,
		jira:      jira,
		store:     store,
	}
	go func() {
		checkAcks(s)
	}()
	return s
}

func checkAcks(s *service) {
	for true {
		time.Sleep(10 * time.Second)
		acks, err := s.store.getAcks()
		if err != nil {
			s.alert.SendError(context.Background(), err)
			continue
		}
		for _, ack := range acks {
			if time.Since(ack.LastSent) > (2 * time.Minute) {
				if ack.Count == 3 {
					s.alert.SendAlert(context.Background(), ack.Chat, "I have tried callout 3 times and have not received an acknowledgement... For now, I am giving up and going home.")
					err = s.store.DeleteAck(ack.Chat)
					if err != nil {
						s.alert.SendError(context.Background(), err)
					}

					continue
				}
				s.alert.SendAlert(context.Background(), ack.Chat, "Callout has not been acknowledged. I am going to phone again.")
				s.alexa.ResetLastCall(ack.Chat)
				err = s.alexa.SendAlert(context.Background(), ack.Chat, ack.Phone, ack.Name, ack.Fields)
				err = s.store.Bump(ack.Chat)
			}
		}
	}
}

// InvokeCallout will invoke snmp if configured, then create a jira ticket if configured, finally it will invoke a phone
// call via alexa connect, if configured.
func (s *service) InvokeCallout(ctx context.Context, chat uint32, title, message string, ack bool) error {
	s.alert.SendAlert(ctx, chat, fmt.Sprintf("invoking callout for: %s, %s", title, message))
	if s.snmp != nil {
		s.snmp.SendSNMPMessage(ctx, chat)
	}
	name, phone, err := s.firstcall.GetFirstCall(ctx, chat)
	if err != nil {
		s.alert.SendError(ctx, err)
		name = "DEFAULT"
	}
	if s.jira != nil {
		s.jira.CreateJira(ctx, chat, title, message, name)
	}
	if s.alexa != nil {
		m := map[string]string{}
		m["Message"] = message
		err := s.alexa.SendAlert(ctx, chat, phone, name, m)
		if err == nil && ack {
			s.alert.SendAlert(ctx, chat, "Please acknowledge the callout with /ack.")
			err = s.store.AddAck(m, chat, phone, name)
			if err != nil {
				s.alert.SendError(ctx, err)
			}
		}
		return err
	}
	return nil
}
