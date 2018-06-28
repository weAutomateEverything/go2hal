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
	alertMock.EXPECT().SendAlert(context.TODO(), uint32(12345), "resolved\n*alertname*: InstanceDown\n*instance*: armonitor.cloudy.standardbank.co.za:8002\n*job*: armonitor\n*severity*: critical\n*description*: pod armonitor.cloudy.standardbank.co.za:8002 is not available\n")
	service := NewService(alertMock)

	err := service.sendPrometheusAlert(12345, "{  \"receiver\": \"hal\",  \"status\": \"resolved\",  \"alerts\": [    {      \"status\": \"resolved\",      \"labels\": {        \"alertname\": \"InstanceDown\",        \"instance\": \"armonitor.cloudy.standardbank.co.za:8002\",        \"job\": \"armonitor\",        \"severity\": \"critical\"      },      \"annotations\": {        \"description\": \"pod armonitor.cloudy.standardbank.co.za:8002 is not available\"      },      \"startsAt\": \"2018-06-28T06:45:08.297934295Z\",      \"endsAt\": \"2018-06-28T06:50:08.297919919Z\",      \"generatorURL\": \"http://prometheus-5898cc7d9b-4r966:9090/graph?g0.expr=up+%3D%3D+0\u0026g0.tab=1\"    }  ],  \"groupLabels\": {    \"alertname\": \"InstanceDown\"  },  \"commonLabels\": {    \"alertname\": \"InstanceDown\",    \"instance\": \"armonitor.cloudy.standardbank.co.za:8002\",    \"job\": \"armonitor\",    \"severity\": \"critical\"  },  \"commonAnnotations\": {    \"description\": \"pod armonitor.cloudy.standardbank.co.za:8002 is not available\"  },  \"externalURL\": \"http://alertmanager-689788b787-gj7lt:9093\",  \"version\": \"4\",  \"groupKey\": \"{}:{alertname=\\\"InstanceDown\\\"}\"}")

	if err != nil {
		log.Println(err)
		t.Fail()
	}

}
