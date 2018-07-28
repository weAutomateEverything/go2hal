package halaws

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context"
)

func NewXray(service Service) Service {
	return &x{
		service: service,
	}
}

type x struct {
	service Service
}

func (s *x) SendAlert(ctx context.Context, chatId uint32, destination string, name string, variables map[string]string) (err error) {
	return xray.Capture(ctx, "halaws.SendAlert", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chatid", chatId)
		xray.AddMetadata(ctx, "destination", destination)
		xray.AddMetadata(ctx, "name", name)
		xray.AddMetadata(ctx, "variables", variables)
		return s.service.SendAlert(ctx, chatId, destination, name, variables)
	})
}
