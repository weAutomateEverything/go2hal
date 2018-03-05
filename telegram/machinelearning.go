package telegram

import (
	"context"
	"github.com/zamedic/go2hal/machineLearning"
)

type ml struct {
	ml machineLearning.Service
	s  Service
}

func NewMachineLearning(service machineLearning.Service, s Service) Service {
	return &ml{service, s}
}

func (s *ml) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (err error) {
	s.ml.StoreAction(ctx, "TELEGRAM", map[string]interface{}{"type": "SendMessage", "chatId": chatID, "message": message, "messageID": messageID})
	return s.s.SendMessage(ctx, chatID, message, messageID)
}

func (s *ml) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (err error) {
	s.ml.StoreAction(ctx, "TELEGRAM", map[string]interface{}{"type": "SendMessage", "chatId": chatID, "message": message, "messageID": messageID})
	return s.s.SendMessagePlainText(ctx, chatID, message, messageID)
}

func (s *ml) SendImageToGroup(ctx context.Context, image []byte, group int64) error {
	return s.s.SendImageToGroup(ctx, image, group)
}

func (s *ml) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) {
	s.s.SendKeyboard(ctx, buttons, text, chat)
}

func (s *ml) RegisterCommand(command Command) {
	s.s.RegisterCommand(command)
}

func (s *ml) RegisterCommandLet(commandlet Commandlet) {
	s.s.RegisterCommandLet(commandlet)
}
