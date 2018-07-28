package alert

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context"
)

func NewXray(s Service) Service {
	return &sXray{
		s,
	}
}

type sXray struct {
	Service
}

func (s sXray) SendAlert(ctx context.Context, chatId uint32, message string) (err error) {
	return xray.Capture(ctx, "alert.SendAlert", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chat", chatId)
		xray.AddMetadata(ctx, "message", message)
		return s.Service.SendAlert(ctx, chatId, message)
	})

}

func (s sXray) SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) (err error) {
	return xray.Capture(ctx, "alert.SendImageToAlertGroup", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chatid", chatid)
		return s.Service.SendImageToAlertGroup(ctx, chatid, image)
	})
}

func (s sXray) SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) (err error) {
	return xray.Capture(ctx, "alert.SendDocumentToAlertGroup", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "chat", chatid)
		xray.AddMetadata(ctx, "extension", extension)
		return s.Service.SendDocumentToAlertGroup(ctx, chatid, document, extension)
	})
}

func (s sXray) SendError(ctx context.Context, err error) (errout error) {
	return xray.Capture(ctx, "alert.SendError", func(ctx context.Context) error {
		xray.AddMetadata(ctx, "err", err)
		return s.Service.SendError(ctx, err)

	})
}

func (s sXray) SendErrorImage(ctx context.Context, image []byte) (err error) {
	return xray.Capture(ctx, "alert.SendErrorImage", func(ctx context.Context) error {
		return s.Service.SendErrorImage(ctx, image)
	})
}
