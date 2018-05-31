package github

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"net/http"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)
	sendAlertHandler := kithttp.NewServer(MakeSendAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	r := mux.NewRouter()

	r.Handle("/api/github/{chatid:[0-9]+}/event", sendAlertHandler).Methods("POST")

	return r

}
