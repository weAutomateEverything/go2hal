package callout

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/halmock"
	"github.com/weAutomateEverything/go2hal/jira"
	snmp2 "github.com/weAutomateEverything/go2hal/snmp"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestService_FirstCall(t *testing.T) {
	testData, err := ioutil.ReadFile("testData/viewcallout.asp.html")
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testData)
	}))
	defer ts.Close()

	os.Setenv("CALLOUT_URL", ts.URL)

	ctrl := gomock.NewController(t)
	alert := alert.NewMockService(ctrl)
	snmp := snmp2.NewMockService(ctrl)
	jira := jira.NewMockService(ctrl)
	aws := halaws.NewMockService(ctrl)

	alert.EXPECT().SendError(context.TODO(), halmock.ErrorMsgMatches(errors.New("invoking callout for: Test, Sample")))
	snmp.EXPECT().SendSNMPMessage(context.TODO())
	jira.EXPECT().CreateJira(context.TODO(), "Test", "Sample", "BOB1")
	aws.EXPECT().SendAlert("+27841231234")

	svc := NewService(alert, snmp, jira, aws)
	svc.InvokeCallout(context.TODO(), "Test", "Sample")

}
