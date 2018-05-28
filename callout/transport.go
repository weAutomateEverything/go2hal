package callout

import (
	"context"
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	calloutHandler := kithttp.NewServer(makeCalloutEndpoint(service), decodeCalloutRequest, gokit.EncodeResponse, opts...)
	r := mux.NewRouter()

	r.Handle("/callout/", calloutHandler).Methods("POST")

	return r
}

func decodeCalloutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	v := &SendCalloutRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
