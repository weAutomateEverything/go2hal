package firstCall

import (
	"context"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
)

type Service interface {
	GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error)
	AddCalloutFunc(function CalloutFunction)
	IsConfigured(chat uint32) bool
	Escalate(ctx context.Context, chat uint32, count int) (name string, number string, err error)
}

type CalloutFunction interface {
	GetFirstCallDetails(ctx context.Context, chat uint32) (name string, number string, err error)
	Escalate(ctx context.Context, count int, chat uint32) (name string, number string, err error)
	Configured(chat uint32) bool
}

func NewCalloutService() Service {
	v := &calloutService{}
	v.services = make([]CalloutFunction, 0)
	return v
}

type calloutService struct {
	services []CalloutFunction
}

func (s *calloutService) Escalate(ctx context.Context, chat uint32, count int) (name string, number string, err error) {
	for _, callout := range s.services {
		name, number, err = callout.Escalate(ctx, count, chat)
		if err == nil {
			return
		}
	}
	return
}

func (s *calloutService) IsConfigured(chat uint32) bool {
	for _, callout := range s.services {
		if callout.Configured(chat) {
			return true
		}
	}
	return false
}

func (s *calloutService) GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error) {
	for _, callout := range s.services {
		name, number, err = callout.GetFirstCallDetails(ctx, chat)
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("no callout has been defined for %v", chat)
	return
}

func (s *calloutService) AddCalloutFunc(function CalloutFunction) {
	s.services = append([]CalloutFunction{function}, s.services...)
}

//-----
type DefaultCalloutService interface {
	setDefaultCallout(ctx context.Context, chat uint32, number string) error
}

type defaultFirstCallService struct {
	store Store
	alert alert.Service
}

func (s *defaultFirstCallService) Escalate(ctx context.Context, count int, chat uint32) (name string, number string, err error) {
	if count > 3 {
		err = errors.New("No escalation defined for default callout. Giving up. ")
		return
	}
	return s.GetFirstCallDetails(ctx, chat)

}

func (s *defaultFirstCallService) Configured(chat uint32) bool {
	_, err := s.store.getDefaultNumber(chat)
	return err == nil
}

func (s *defaultFirstCallService) setDefaultCallout(ctx context.Context, chat uint32, number string) (err error) {
	err = s.store.setDefaultNumber(chat, number)
	if err != nil {
		return
	}
	s.alert.SendAlert(ctx, chat, fmt.Sprintf("Default Callout for your group has been set to %v", number))
	return
}

func NewDefaultFirstcallService(store Store, service alert.Service) CalloutFunction {
	return &defaultFirstCallService{
		store: store,
		alert: service,
	}
}

func (s *defaultFirstCallService) GetFirstCallDetails(ctx context.Context, chat uint32) (name string, number string, err error) {
	number, err = s.store.getDefaultNumber(chat)
	name = "DEFAULT"
	return
}
