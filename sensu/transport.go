package sensu

import (
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
	"context"
	"encoding/json"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(gokit.EncodeError),
	}

	sensuAlert := kithttp.NewServer(makeSensuEndpoint(service), decodeSensu, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/sensu", sensuAlert).Methods("POST")

	return r

}


func decodeSensu(_ context.Context, r *http.Request) (interface{}, error) {
	var sensu SensuMessageRequest
	err := json.NewDecoder(r.Body).Decode(&sensu)
	return sensu, err
}
