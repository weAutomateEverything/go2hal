package appdynamics

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/ssh"
	"golang.org/x/net/context"
	"testing"
)

func TestServiceWithQuotes(t *testing.T) {

	ctrl := gomock.NewController(t)

	alertService := alert.NewMockService(ctrl)
	sshService := ssh.NewMockService(ctrl)
	store := NewMockStore(ctrl)
	calloutService := callout.NewMockService(ctrl)

	service := NewService(alertService, sshService, store, calloutService)

	store.EXPECT().getAllEndpoints().AnyTimes().Return(nil, nil)
	alertService.EXPECT().SendAlert(context.TODO(), uint32(12345), "🔴  [Error]  - ServletException: java.lang.NumberFormatException: For input string: \"\"\n*Application:* TESTAPP\n*Tier:* TESTTIER\n*Node:* test\n")

	d := AppdynamicsMessage{
		Events: []Event{
			{
				Severity: "ERROR",
				Application: Name{
					Name: "TESTAPP",
				},
				Tier: Name{
					Name: "TESTTIER",
				},
				Node: Name{
					Name: "test",
				},
				DisplayName:  "Business Transaction Error",
				EventMessage: "[Error]  - ServletException: java.lang.NumberFormatException: For input string: \"\"",
			},
		},
	}

	err := service.sendAppdynamicsAlert(context.TODO(), uint32(12345), d)

	if err != nil {
		t.Error(err)
	}
}
