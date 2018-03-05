package alert

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/zamedic/go2hal/telegram"
	"golang.org/x/net/context"
	"testing"
)

func TestService_SendAlert(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockStore := NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockStore)
	ctx := context.TODO()

	mockStore.EXPECT().alertGroup().Return(int64(12345), nil)
	mockTelegram.EXPECT().SendMessage(ctx, int64(12345), "Hello World", 0)

	svc.SendAlert(ctx, "Hello World")
}

func TestService_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockStore := NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockStore)
	ctx := context.TODO()

	mockStore.EXPECT().heartbeatGroup().Return(int64(12345), nil)
	mockTelegram.EXPECT().SendMessagePlainText(ctx, int64(12345), "💩  Hello World", 0)

	svc.SendError(ctx, errors.New("Hello World"))
}
