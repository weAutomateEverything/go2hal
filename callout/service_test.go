package callout

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/jira"
	snmp2 "github.com/weAutomateEverything/go2hal/snmp"
	"golang.org/x/net/context"
	"testing"
)

func TestService_FirstCall(t *testing.T) {

	ctrl := gomock.NewController(t)
	alert := alert.NewMockService(ctrl)
	snmp := snmp2.NewMockService(ctrl)
	jira := jira.NewMockService(ctrl)
	aws := halaws.NewMockService(ctrl)
	firstCallService := firstCall.NewMockService(ctrl)
	store := NewMockStore(ctrl)

	alert.EXPECT().SendAlert(context.TODO(), uint32(12345), "invoking callout for: Test, Sample")
	alert.EXPECT().SendAlert(context.TODO(), uint32(12345), "Please acknowledge the callout with /ack.")
	snmp.EXPECT().SendSNMPMessage(context.TODO(), uint32(12345))
	jira.EXPECT().CreateJira(context.TODO(), uint32(12345), "Test", "Sample", "BOB1")
	aws.EXPECT().SendAlert(context.TODO(), uint32(12345), "+27841231234", "BOB1", map[string]string{"Message": "Sample"})
	firstCallService.EXPECT().GetFirstCall(context.TODO(), uint32(12345)).Return("BOB1", "+27841231234", nil)
	store.EXPECT().AddAck(map[string]string{"Message": "Sample"}, uint32(12345), "+27841231234", "BOB1")

	svc := NewService(alert, firstCallService, snmp, jira, aws, store)
	svc.InvokeCallout(context.TODO(), uint32(12345), "Test", "Sample", true)

}
