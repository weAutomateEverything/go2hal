package skynet

import (
	"context"
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
	"github.com/zamedic/go2hal/machineLearning"
	"net/http"
)

//MakeHandler returns a rest http handler.
//machine learning service can be set to nil if you do not want to store the requests
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	skynetRebuild := kithttp.NewServer(makeSkynetRebuildEndpoint(service), decodeSkynetRebuildRequest, gokit.EncodeResponse, opts...)
	skynetAlert := kithttp.NewServer(makeSkynetAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/skynet/", skynetAlert).Methods("POST")
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
