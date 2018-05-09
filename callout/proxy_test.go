package callout

import (
	"golang.org/x/net/context"
	"os"
)

func ExampleNewCalloutProxy() {
	os.Setenv("HAL_ENDPOINT", "http://<some server>:8000")
	svc := NewCalloutProxy()
	svc.InvokeCallout(context.TODO(), "Test Message", "The body of the message goes in here", nil)
}

func ExampleNewKubernetesCalloutProxy() {
	svc := NewKubernetesCalloutProxy("")
	svc.InvokeCallout(context.TODO(), "Test Message", "The body of the message goes in here", nil)
}
