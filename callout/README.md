# HAL Callout Service

The service allows for callout to be invoked via 4 different methods.
1. Telegram Message
1. JIRA Task
1. SNMP Track
1. Phone Call via Amazon AWS Connect

## Server
### Service
If you wish to include the callout service within your code, you would need to create a new instance of the service. 

```go
	calloutService := callout.NewService(alertService, snmpService, jiraService, aws)
```
snmp, jira and aws can be set to nil if you do not wish to use them. 

### REST Server
A rest server is available, to configure 
```go
mux := http.NewServeMux()
mux.Handle("/callout/", callout.MakeHandler(calloutService, httpLogger, machineLearningService))
```

machine learning service can be set to nil if you do not wish to record the requests.

#### Invoke Callout
Rest endpoint is available at `http://<hal server>:8000/callout/`.
Send a POST request with JSON Format
```json
{
    "Title" : "some title",
    "Message": "some message"
}
```


##Client
A remote HTTP proxy client is provided for easy access to the callout.

Examples are located withint he proxy_test.go file

### Standard HTTP
```go
func ExampleNewCalloutProxy() {
	os.Setenv("HAL_ENDPOINT","http://<some server>:8000")
	svc := NewCalloutProxy()
	svc.InvokeCallout(context.TODO(),"Test Message","The body of the message goes in here")
}
```

### Kubernetes 
```go
func ExampleNewKubernetesCalloutProxy() {
	svc := NewKubernetesCalloutProxy("")
	svc.InvokeCallout(context.TODO(),"Test Message","The body of the message goes in here")
}

```
