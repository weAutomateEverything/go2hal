package halaws

import (
	"github.com/weAutomateEverything/go2hal/machineLearning"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/gorilla/mux"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)
	sendAlertHandler := kithttp.NewServer(MakeSendAlertEndpoint(service),gokit.DecodeString,gokit.EncodeResponse,opts...)
	r := mux.NewRouter()

	r.Handle("/aws/sendTestAlert",sendAlertHandler).Methods("POST")

	return r

}