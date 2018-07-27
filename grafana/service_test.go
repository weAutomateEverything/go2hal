package grafana

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"log"
	"testing"
)

func TestService_Grafana(t *testing.T) {
	ctrl := gomock.NewController(t)

	alertMock := alert.NewMockService(ctrl)
	alertMock.EXPECT().SendAlert(context.TODO(), uint32(12345), "*Alerting* Test notification\nSomeone is testing the alert notification within grafana.")

	service := NewService(alertMock)

	err := service.sendGrafanaAlert(context.TODO(), 12345, "{\"evalMatches\":[{\"value\":100,\"metric\":\"High value\",\"tags\":null},{\"value\":200,\"metric\":\"Higher Value\",\"tags\":null}],\"imageUrl\":\"http://grafana.org/assets/img/blog/mixed_styles.png\",\"message\":\"Someone is testing the alert notification within grafana.\",\"ruleId\":0,\"ruleName\":\"Test notification\",\"ruleUrl\":\"http://localhost:3000/\",\"state\":\"alerting\",\"title\":\"[Alerting] Test notification\"}")

	if err != nil {
		log.Println(err)
		t.Fail()
	}

}
