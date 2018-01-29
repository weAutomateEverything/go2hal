package callout

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/telegram"
)

type whosOnFirstCall struct {
	alert alert.Service
	telegram telegram.Service
	service Service

}

func NewWhosOnFirstCallCommand(alert alert.Service,telegram telegram.Service,service Service) telegram.Command{
	return &whosOnFirstCall{alert,telegram,service}
}


/* Set Heartbeat group */
func (s *whosOnFirstCall) CommandIdentifier() string {
	return "FirstCall"
}

func (s *whosOnFirstCall) CommandDescription() string {
	return "Who is on first call?"
}

func (s *whosOnFirstCall) Execute(update tgbotapi.Update) {
	name, err := s.service.getFirstCallName()
	if err != nil {
		s.alert.SendError(err)
		return
	}
	s.telegram.SendMessage(update.Message.Chat.ID, fmt.Sprintf("%s is on first call",name), update.Message.MessageID)
}