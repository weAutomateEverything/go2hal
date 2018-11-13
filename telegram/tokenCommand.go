package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
)

type token struct {
	Service
	Store
}

func (token) CommandIdentifier() string {
	return "token"
}

func (token) CommandDescription() string {
	return "get a token to perform config operations via the API"
}

func (token) RestrictToAuthorised() bool {
	return true
}

func (token) Show(chat uint32) bool {
	return true
}

func (s token) Execute(ctx context.Context, update tgbotapi.Update) {
	chat, err := s.GetUUID(update.Message.Chat.ID, update.Message.Chat.Title)
	if err != nil {
		s.SendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf("we had an error fetching your room details: %v", err.Error()), update.Message.MessageID)
		return
	}
	token, err := makeToken(chat)
	if err != nil {
		s.SendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf("we had an error generating your token: %v", err.Error()), update.Message.MessageID)
		return
	}
	s.SendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf("*JWT Token*\n"+
		"To use the token, add it to the HTTP Request header as \n"+
		"Authorization: Bearer <token>\n\n"+
		"%v", token), update.Message.MessageID)
}

func NewTokenCommand(service Service, store Store) Command {
	return &token{
		service, store,
	}
}
