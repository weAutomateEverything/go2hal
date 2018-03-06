package examples

import (
	"github.com/weAutomateEverything/go2hal/alert"
	"os"
)

func clientStandalone() {
	// This can be set a number of ways, below is just an example of setting en environment variable
	os.Setenv("ALERT_ENDPOINT", "http://localhost:8000/")
	alertProxy := alert.NewAlertProxy()
	alertProxy.SendAlert(nil, "Hello World")
}

func clientKubernetes() {
	//Leaving the namespace blank, means use the same namespace of the current pod.
	alertProxy := alert.NewKubernetesAlertProxy("")
	alertProxy.SendAlert(nil, "Hello World")
}
