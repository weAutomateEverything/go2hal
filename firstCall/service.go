package firstCall

import (
	"context"
	"os"
)

type Service interface {
	GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error)
}

type defaultFirstCallService struct {
}

func NewDefaultFirstcallService() Service {
	return &defaultFirstCallService{}
}

func (*defaultFirstCallService) GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error) {
	return "DEFAULT", os.Getenv("DEFAULT_CALLOUT_NUMBER"), nil
}
