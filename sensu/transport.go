package sensu

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)

	endpoint := makeSensuEndpoint(service)

	sensuAlert := kithttp.NewServer(endpoint, decodeSensu, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/sensu/", sensuAlert).Methods("POST")

	return r

}

func decodeSensu(_ context.Context, r *http.Request) (interface{}, error) {
	var sensu SensuMessageRequest
	err := json.NewDecoder(r.Body).Decode(&sensu)
	return sensu, err
}
