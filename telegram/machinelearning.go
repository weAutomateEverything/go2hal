package telegram

import (
	"context"
	"github.com/weAutomateEverything/go2hal/machineLearning"
)

//NewMachineLearning returns a decorated Service object that will log the telegram actions executed to the machine
//learning database
func NewMachineLearning(service machineLearning.Service, s Service) Service {
	return &ml{service, s}
}

type ml struct {
	ml machineLearning.Service
	s  Service
}

func (s *ml) SendMessageWithCorrelation(ctx context.Context, chatID int64, message string, messageID int, correlationId string) (msgid int, err error) {
	return s.s.SendMessageWithCorrelation(ctx, chatID, message, messageID, correlationId)
}

func (s *ml) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	s.ml.StoreAction(ctx, "TELEGRAM", map[string]interface{}{"type": "SendMessage", "chatId": chatID, "message": message, "messageID": messageID})
	return s.s.SendMessage(ctx, chatID, message, messageID)
}

func (s *ml) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	s.ml.StoreAction(ctx, "TELEGRAM", map[string]interface{}{"type": "SendMessage", "chatId": chatID, "message": message, "messageID": messageID})
	return s.s.SendMessagePlainText(ctx, chatID, message, messageID)
}

func (s *ml) SendImageToGroup(ctx context.Context, image []byte, group int64) error {
	return s.s.SendImageToGroup(ctx, image, group)
}

func (s *ml) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (int, error) {
	return s.s.SendKeyboard(ctx, buttons, text, chat)
}

func (s *ml) SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) error {
	return s.s.SendDocumentToGroup(ctx, document, extension, group)
}

func (s *ml) RegisterCommand(command Command) {
	s.s.RegisterCommand(command)
}

func (s *ml) RegisterCommandLet(commandlet Commandlet) {
	s.s.RegisterCommandLet(commandlet)
}

func (s *ml) requestAuthorisation(ctx context.Context, chat uint32, name string) (string, error) {
	return s.s.requestAuthorisation(ctx, chat, name)
}

func (s *ml) pollAuthorisation(token string) (uint32, error) {
	return s.s.pollAuthorisation(token)
}
