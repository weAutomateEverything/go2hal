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

}


/* Set Heartbeat group */
func (s *whosOnFirstCall) commandIdentifier() string {
	return "FirstCall"
}

func (s *whosOnFirstCall) commandDescription() string {
	return "Who is on first call?"
}

func (s *whosOnFirstCall) execute(update tgbotapi.Update) {
	name, err := getFirstCallName()
	if err != nil {
		s.alert.SendError(err)
		return
	}
	s.telegram.SendMessage(update.Message.Chat.ID, fmt.Sprintf("%s is on first call",name), update.Message.MessageID)
}