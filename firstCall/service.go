package firstCall

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
)

type Service interface {
	GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error)
	AddCalloutFunc(function CalloutFunction)
}

type CalloutFunction interface {
	GetFirstCallDetails(ctx context.Context, chat uint32) (name string, number string, err error)
}

func NewCalloutService() Service {
	v := &calloutService{}
	v.services = []CalloutFunction{newDefaultFirstcallService()}
	return v
}

type calloutService struct {
	services []CalloutFunction
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

type defaultFirstCallService struct {
}

func newDefaultFirstcallService() CalloutFunction {
	return &defaultFirstCallService{}
}

func (*defaultFirstCallService) GetFirstCallDetails(ctx context.Context, chat uint32) (name string, number string, err error) {
	num := os.Getenv("DEFAULT_CALLOUT_NUMBER")
	if num == "" {
		err = errors.New("No DEFAULT_CALLOUT_NUMBER environment variable has been defined.")
		return
	}
	return "DEFAULT", os.Getenv("DEFAULT_CALLOUT_NUMBER"), nil
}
