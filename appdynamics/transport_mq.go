package appdynamics

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

//MakeHandler returns a restful http handler for the appdynamics service
//the Machine Learning service can be set to nil if you do not wish to log the http requests
func MakeMqHandler(service MqService, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)
	addAppdynamicsQueueEndpoint := kithttp.NewServer(makeAddAppdynamicsQueueEndpoint(service), decodeAddMqEndpointRequest, gokit.EncodeResponse, opts...)
	r := mux.NewRouter()

	// swagger:operation POST /api/appdynamics/{chatid}/queue appdynamics addQueue
	//
	// Add a IBM MQ queue to monitor
	//
	// Add a ibm MQ queue thats been monitored by the App Dynamics MQ Extension (https://www.appdynamics.com/community/exchange/extension/websphere-mq-monitoring-extension/).
	// Once the plugin has been configured, it should start sending MQ details.
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// - name: body
	//   description: queue details
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/AddAppdynamicsQueueEndpointRequest"
	// responses:
	//   '200':
	//     description: Command has been executed successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/appdynamics/{chatid:[0-9]+}/queue", addAppdynamicsQueueEndpoint).Methods("POST")

	return r
}

func decodeAddMqEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := AddAppdynamicsQueueEndpointRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
