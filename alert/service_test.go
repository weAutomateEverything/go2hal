package alert

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/zamedic/go2hal/telegram"
	"testing"
)

func TestService_SendAlert(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockStore := NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockStore)

	mockStore.EXPECT().alertGroup().Return(int64(12345), nil)
	mockTelegram.EXPECT().SendMessage(int64(12345), "Hello World", 0)

	svc.SendAlert("Hello World")
}

func TestService_SendError(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockTelegram := telegram.NewMockService(ctrl)
	mockStore := NewMockStore(ctrl)

	svc := NewService(mockTelegram, mockStore)

	mockStore.EXPECT().heartbeatGroup().Return(int64(12345), nil)
	mockTelegram.EXPECT().SendMessagePlainText(int64(12345), "💩  Hello World", 0)

	svc.SendError(errors.New("Hello World"))
}
