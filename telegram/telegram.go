package telegram

import (
	"fmt"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v4"
)

func NewTelegramAuthApprovalCommand(service2 Service, store Store) Commandlet {
	return &approveAuth{
		service2,
		store,
	}
}

type approveAuth struct {
	Service
	Store
}

func (s approveAuth) CanExecute(update tgbotapi.Update, state State) bool {
	return update.Message.Text == "Approve access" && update.Message.ReplyToMessage != nil
}

func (s approveAuth) Execute(update tgbotapi.Update, state State) {
	err := s.approveAuthRequest(update.Message.ReplyToMessage.MessageID, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.ID)
	if err != nil {
		s.SendMessagePlainText(context.TODO(), update.Message.Chat.ID,
			fmt.Sprintf("There was an error approving your request. %v", err.Error()), update.Message.MessageID)

	} else {
		s.SendMessagePlainText(context.TODO(), update.Message.Chat.ID,
			"The access request was successfully approved", update.Message.MessageID)
	}
}

func (s approveAuth) NextState(update tgbotapi.Update, state State) string {
	return ""
}

func (s approveAuth) Fields(update tgbotapi.Update, state State) []string {
	return nil
}
