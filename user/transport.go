package user

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
	"github.com/zamedic/go2hal/machineLearning"
	"net/http"
)

/*
MakeHandler returns a HTTP Restul endpoint to handle user requests
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	addBulkUserEndpoint := kithttp.NewServer(makeBulkUserUploadEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/users/", addBulkUserEndpoint).Methods("POST")

	return r

}
