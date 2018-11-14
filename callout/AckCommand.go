package callout

import (
	"context"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/telegram"
	"gopkg.in/telegram-bot-api.v4"
)

func NewAckCommand(alert alert.Service, telegram telegram.Store, store Store) telegram.Command {
	return ackCommand{
		alert:    alert,
		telegram: telegram,
		store:    store,
	}
}

type ackCommand struct {
	alert    alert.Service
	telegram telegram.Store
	store    Store
}

func (ackCommand) CommandIdentifier() string {
	return "ack"
}

func (ackCommand) CommandDescription() string {
	return "Acknowledge a callout"
}

func (ackCommand) RestrictToAuthorised() bool {
	return true
}

func (ackCommand) Show(chat uint32) bool {
	return true
}

func (s ackCommand) Execute(ctx context.Context, update tgbotapi.Update) {
	uuid, err := s.telegram.GetUUID(update.Message.Chat.ID, update.Message.Chat.Title)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}
	err = s.store.DeleteAck(uuid)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}
	s.alert.SendAlert(ctx, uuid, "Thank you. Callout has been successfully acknowledged.")
}
