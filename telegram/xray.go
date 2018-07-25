package telegram

import (
	"context"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func NewXray(s Service) Service {
	return sXray{
		s,
	}
}

type sXray struct {
	Service
}

func (s sXray) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendMessage")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendMessage(ctx, chatID, message, messageID)
}
func (s sXray) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendMessagePlainText")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendMessagePlainText(ctx, chatID, message, messageID)
}
func (s sXray) SendImageToGroup(ctx context.Context, image []byte, group int64) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendImageToGroup")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendImageToGroup(ctx, image, group)
}
func (s sXray) SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendDocumentToGroup")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendDocumentToGroup(ctx, document, extension, group)
}

func (s sXray) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (message int, err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendKeyboard")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendKeyboard(ctx, buttons, text, chat)
}
func (s sXray) RegisterCommand(command Command) {
	s.Service.RegisterCommand(command)
}
func (s sXray) RegisterCommandLet(commandlet Commandlet) {
	s.Service.RegisterCommandLet(commandlet)
}

func (s sXray) requestAuthorisation(chat uint32, name string) (string, error) {
	return s.Service.requestAuthorisation(chat, name)
}
func (s sXray) pollAuthorisation(token string) (uint32, error) {
	return s.Service.pollAuthorisation(token)
}
