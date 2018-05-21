package alert

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"os"
	"testing"
)

func TestService_SendAlert(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockTelegramStore := telegram.NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockTelegramStore)
	ctx := context.TODO()

	mockTelegram.EXPECT().SendMessage(ctx, int64(12345), "Hello World", 0)
	mockTelegramStore.EXPECT().GetRoomKey(uint32(54321)).Return(int64(12345), nil)

	svc.SendAlert(ctx, uint32(54321), "Hello World")
}

func TestService_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockStore := telegram.NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockStore)
	ctx := context.TODO()

	os.Setenv("ERROR_GROUP", "12345")

	mockTelegram.EXPECT().SendMessagePlainText(ctx, int64(12345), "💩  Hello World", 0)

	svc.SendError(ctx, errors.New("Hello World"))
}
