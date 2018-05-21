package github

import (
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
)

type service struct {
	alert alert.Service
}

func NewService(alert alert.Service) Service {
	return &service{
		alert: alert,
	}
}

type Service interface {
	sendGithubMessage(ctx context.Context, chatId uint32, message string)
}

func (s *service) sendGithubMessage(ctx context.Context, chatId uint32, message string) {
	s.alert.SendAlert(ctx, chatId, message)
}
