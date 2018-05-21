package sensu

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
)

//MakeHandler retuns a http rest request handler for sensu
//the machine learning service can be nil if you do not wish to save the request message
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	endpoint := makeSensuEndpoint(service)

	sensuAlert := kithttp.NewServer(endpoint, decodeSensu, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/sensu/{chatid:[0-9]+}", sensuAlert).Methods("POST")

	return r

}

func decodeSensu(_ context.Context, r *http.Request) (interface{}, error) {
	var sensu SensuMessageRequest
	err := json.NewDecoder(r.Body).Decode(&sensu)
	return sensu, err
}
