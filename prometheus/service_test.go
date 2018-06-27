package prometheus

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
	alertMock.EXPECT().SendAlert(context.TODO(), uint32(12345), "resolved\nsome label\nsome annotation")

	service := NewService(alertMock)

	err := service.sendPrometheusAlert(12345, "{  \"version\": \"4\",  \"groupKey\": \"test\",  \"status\": \"<resolved|firing>\",  \"receiver\": \"test\",  \"groupLabels\": \"test\",  \"commonLabels\": \"test\",  \"commonAnnotations\": \"test\",  \"externalURL\":\"test\",  \"alerts\": [    {      \"status\": \"resolved\",      \"labels\": \"some label\",      \"annotations\": \"some annotation\",      \"startsAt\": \"<rfc3339>\",      \"endsAt\": \"<rfc3339>\",      \"generatorURL\": \"some url\"    }  ]}")

	if err != nil {
		log.Println(err)
		t.Fail()
	}

}
