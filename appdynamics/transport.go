package appdynamics

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

/*
MakeHandler returns a restful http handler for the appdynamics service

the Machine Learning service can be set to nil if you do not wish to log the http requests
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	appDynamicsAlertEndoint := kithttp.NewServer(makeAppDynamicsAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	addAppdynamicsEndpoint := kithttp.NewServer(makeAddAppdynamicsEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	addAppdynamicsQueueEndpoint := kithttp.NewServer(makeAddAppdynamicsQueueEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	executeCommandFromAppdynamics := kithttp.NewServer(makExecuteCommandFromAppdynamics(service), decodeExecuteRequest, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/appdynamics/", appDynamicsAlertEndoint).Methods("POST")
	r.Handle("/appdynamics/system/queue", addAppdynamicsEndpoint).Methods("POST")
	r.Handle("/appdynamics/system", addAppdynamicsQueueEndpoint).Methods("POST")
	r.Handle("/appdynamics/execute", executeCommandFromAppdynamics).Methods("POST")

	return r

}

func decodeExecuteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := ExecuteAppDynamicsCommandRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
