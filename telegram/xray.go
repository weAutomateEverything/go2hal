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
	xray.Capture(ctx, "telegram.SendMessage", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chat", chatID)
		xray.AddMetadata(ctx, "message", message)
		xray.AddMetadata(ctx, "messageid", messageID)
		msgid, err = s.Service.SendMessage(ctx, chatID, message, messageID)
		return err
	})
	return
}
func (s sXray) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	xray.Capture(ctx, "telegram.SendMessagePlainText", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chatid", chatID)
		xray.AddMetadata(ctx, "message", message)
		xray.AddMetadata(ctx, "messageid", messageID)
		msgid, err = s.Service.SendMessagePlainText(ctx, chatID, message, messageID)
		return err
	})
	return
}
func (s sXray) SendImageToGroup(ctx context.Context, image []byte, group int64) (err error) {
	return xray.Capture(ctx, "telegram.SendImageToGroup", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "group", group)
		return s.Service.SendImageToGroup(ctx, image, group)

	})

}
func (s sXray) SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) (err error) {
	return xray.Capture(ctx, "telegram.SendDocumentToGroup", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "extension", extension)
		xray.AddMetadata(ctx, "group", group)
		return s.Service.SendDocumentToGroup(ctx, document, extension, group)
	})

}

func (s sXray) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (message int, err error) {
	xray.Capture(ctx, "telegram.SendKeyboard", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "buttons", buttons)
		xray.AddMetadata(ctx, "text", text)
		xray.AddMetadata(ctx, "chat", chat)
		message, err = s.Service.SendKeyboard(ctx, buttons, text, chat)
		return err
	})
	return

}
func (s sXray) RegisterCommand(command Command) {
	s.Service.RegisterCommand(command)
}
func (s sXray) RegisterCommandLet(commandlet Commandlet) {
	s.Service.RegisterCommandLet(commandlet)
}

func (s sXray) requestAuthorisation(ctx context.Context, chat uint32, name string) (string, error) {
	return s.Service.requestAuthorisation(ctx, chat, name)
}
func (s sXray) pollAuthorisation(token string) (uint32, error) {
	return s.Service.pollAuthorisation(token)
}
