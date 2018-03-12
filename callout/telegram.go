package callout

import (
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v4"
)

type whosOnFirstCall struct {
	alert    alert.Service
	telegram telegram.Service
	service  Service
}

func NewWhosOnFirstCallCommand(alert alert.Service, telegram telegram.Service, service Service) telegram.Command {
	return &whosOnFirstCall{alert, telegram, service}
}

/* Set Heartbeat group */
func (s *whosOnFirstCall) CommandIdentifier() string {
	return "FirstCall"
}

func (s *whosOnFirstCall) CommandDescription() string {
	return "Who is on first call?"
}

func (s *whosOnFirstCall) Execute(update tgbotapi.Update) {
	name, phone, err := s.service.getFirstCall(context.TODO())
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return
	}
	s.telegram.SendMessage(context.TODO(), update.Message.Chat.ID, fmt.Sprintf("%s is on first call. Number %v", name,phone), update.Message.MessageID)
}
