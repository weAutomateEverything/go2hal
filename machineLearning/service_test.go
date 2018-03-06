package machineLearning

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestService_StoreAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := NewMockStore(ctrl)

	svc := NewService(store)
	ctx := context.WithValue(context.TODO(), keyPrincipalID, "12345")
	m := map[string]interface{}{"test1": 1}

	store.EXPECT().SaveAction("12345", "TEST", gomock.Any(), m)

	svc.StoreAction(ctx, "TEST", m)
}

func TestService_StoreHTTPRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := NewMockStore(ctrl)

	svc := NewService(store)
	ctx := context.TODO()
	req := http.Request{Method: "POST", Body: ioutil.NopCloser(bytes.NewBuffer([]byte("Hello World")))}
	store.EXPECT().SaveInputRecord("HTTP", gomock.Any(), map[string]interface{}{"body": "Hello World", "method": "POST"}).Return("KEY")

	ctx = svc.StoreHTTPRequest(ctx, &req)

	if ctx.Value(keyPrincipalID) != "KEY" {
		t.Errorf("Invalid Value found in context. Expected 'KEY', received %v", ctx.Value(keyPrincipalID))
	}

}
