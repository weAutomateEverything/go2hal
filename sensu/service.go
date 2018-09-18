package sensu

import (
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"strings"
)

type Service interface {
	handleSensu(ctx context.Context, chatId uint32, sensu SensuMessageRequest)
}

type service struct {
	alert alert.Service
}

func NewService(alert alert.Service) Service {
	return &service{alert: alert}
}

func (s *service) handleSensu(ctx context.Context, chatId uint32, sensu SensuMessageRequest) {
	for _, msg := range sensu.Attachments {
		e := ""
		if strings.Index(msg.Title, "CRITICAL") > 0 {
			e = ":rotating_light:"

		} else if strings.Index(msg.Title, "WARNING") > 0 {
			e = ":warning:"
		} else {
			e = ":white_check_mark:"
		}

		msg.Title = strings.Replace(msg.Title, "*", "\\*", -1)
		msg.Title = strings.Replace(msg.Title, "_", "\\_", -1)

		msg.Text = strings.Replace(msg.Text, "*", "\\*", -1)
		msg.Text = strings.Replace(msg.Text, "_", "\\_", -1)

		m := emoji.Sprintf("%v *%v*\n %v", e, msg.Title, msg.Text)
		s.alert.SendAlert(ctx, chatId, m)

	}

}
