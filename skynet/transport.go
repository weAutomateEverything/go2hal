package skynet

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
	opts := gokit.GetServerOpts(logger)

	skynetRebuild := kithttp.NewServer(makeSkynetRebuildEndpoint(service), decodeSkynetRebuildRequest, gokit.EncodeResponse, opts...)
	skynetAlert := kithttp.NewServer(makeSkynetAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/skynet", skynetAlert).Methods("POST")
	r.Handle("/skynet/rebuild", skynetRebuild).Methods("POST")

	return r

}

func decodeSkynetRebuildRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SkynetRebuildRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil

}
