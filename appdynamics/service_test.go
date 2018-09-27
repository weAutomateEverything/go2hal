package appdynamics

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/ssh"
	"golang.org/x/net/context"
	"testing"
)

func TestServiceWithQuotes(t *testing.T) {

	ctrl := gomock.NewController(t)

	alertService := alert.NewMockService(ctrl)
	sshService := ssh.NewMockService(ctrl)
	store := NewMockStore(ctrl)

	service := NewService(alertService, sshService, store)

	store.EXPECT().getAllEndpoints().AnyTimes().Return(nil, nil)
	alertService.EXPECT().SendAlert(context.TODO(), uint32(12345), "🔴  [Error]  - ServletException: java.lang.NumberFormatException: For input string: \"\"\n*Application:* TESTAPP\n*Tier:* TESTTIER\n*Node:* test\n")

	err := service.sendAppdynamicsAlert(context.TODO(), uint32(12345), "{\"events\": [{\"severity\": \"ERROR\",\"application\": {\"name\": \"TESTAPP\"},\"tier\": {\"name\": \"TESTTIER\"},\"node\": {\"name\": \"test\"},\"displayName\": \"Business Transaction Error\",\"eventMessage\": \"[Error]  - ServletException: java.lang.NumberFormatException: For input string: \"\"\"}]}")

	if err != nil {
		t.Error(err)
	}
}
