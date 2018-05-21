package alert

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	alertHandler := kithttp.NewServer(makeAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)
	alertErrorHandler := kithttp.NewServer(makeAlertErrorHandler(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/alert/{chatid:[0-9]+}", alertHandler).Methods("POST")
	r.Handle("/alert/{chatid:[0-9]+}/image", imageAlertHandler).Methods("POST")

	r.Handle("/alert/error", alertErrorHandler).Methods("POST")

	return r
}
