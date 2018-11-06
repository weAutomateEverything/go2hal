package firstCall

import (
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v4"
)

type whosOnFirstCall struct {
	alert         alert.Service
	telegram      telegram.Service
	telegramStore telegram.Store
	service       Service
}

func NewWhosOnFirstCallCommand(alert alert.Service, telegram telegram.Service, service Service, store telegram.Store) telegram.Command {
	return &whosOnFirstCall{
		telegramStore: store,
		telegram:      telegram,
		service:       service,
		alert:         alert,
	}
}

func (s *whosOnFirstCall) RestrictToAuthorised() bool {
	return false
}

/* Set Heartbeat group */
func (s *whosOnFirstCall) CommandIdentifier() string {
	return "FirstCall"
}

func (s *whosOnFirstCall) CommandDescription() string {
	return "Who is on first call?"
}

func (s *whosOnFirstCall) Show(chat uint32) bool {
	return s.service.IsConfigured(chat)
}

func (s *whosOnFirstCall) Execute(ctx context.Context, update tgbotapi.Update) {
	uuid, err := s.telegramStore.GetUUID(update.Message.Chat.ID, update.Message.Chat.Title)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}
	name, phone, err := s.service.GetFirstCall(ctx, uuid)
	if err != nil {
		s.alert.SendError(ctx, err)
		s.telegram.SendMessagePlainText(ctx, update.Message.Chat.ID,
			fmt.Sprintf("There was an error fetching your firstcall details. %v", err.Error()),
			update.Message.MessageID)
		return
	}
	s.telegram.SendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf("%s is on first call. Number %v", name, phone), update.Message.MessageID)
}
