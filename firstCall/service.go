package firstCall

import (
	"context"
	"github.com/pkg/errors"
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

func (*defaultFirstCallService) setGroup(ctx context.Context, chat uint32, group string) (name string, number string, err error) {
	return "", "", errors.New("Set Group not implimented for default callout")
}

func (*defaultFirstCallService) getGroup(ctx context.Context, chat uint32) (string, error) {
	return "", errors.New("get Group not implimented for default callout")
}
