package alert

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context"
)

func NewXray(s Service) Service {
	return &s_xray{
		s,
	}
}

type s_xray struct {
	Service
}

func (s s_xray) SendAlert(ctx context.Context, chatId uint32, message string) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendAlert")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendAlert(ctx, chatId, message)
}

func (s s_xray) SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendImageToAlertGroup")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendImageToAlertGroup(ctx, chatid, image)
}

func (s s_xray) SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendDocumentToAlertGroup")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendDocumentToAlertGroup(ctx, chatid, document, extension)
}

func (s s_xray) SendError(ctx context.Context, err error) (errout error) {
	ctx, seg := xray.BeginSegment(ctx, "SendError")
	defer func() {
		seg.Close(errout)
	}()
	return s.Service.SendError(ctx, err)
}

func (s s_xray) SendErrorImage(ctx context.Context, image []byte) (err error) {
	ctx, seg := xray.BeginSegment(ctx, "SendErrorImage")
	defer func() {
		seg.Close(err)
	}()
	return s.Service.SendErrorImage(ctx, image)
}
