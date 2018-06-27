package prometheus

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
)

//MakeHandler retuns a http rest request handler for sensu
//the machine learning service can be nil if you do not wish to save the request message
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	endpoint := makePrometheusAlertEndpoint(service)

	prometheus := kithttp.NewServer(endpoint, gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/api/prometheus/{chatid:[0-9]+}", prometheus).Methods("POST")

	return r

}
