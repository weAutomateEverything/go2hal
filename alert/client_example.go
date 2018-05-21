package alert

import (
	"os"
)

func clientStandalone() {
	// This can be set a number of ways, below is just an example of setting en environment variable
	os.Setenv("ALERT_ENDPOINT", "http://localhost:8000/")
	alertProxy := NewAlertProxy()
	alertProxy.SendAlert(nil, 12345, "Hello World")
}

func clientKubernetes() {
	//Leaving the namespace blank, means use the same namespace of the current pod.
	alertProxy := NewKubernetesAlertProxy("")
	alertProxy.SendAlert(nil, 12345, "Hello World")
}
