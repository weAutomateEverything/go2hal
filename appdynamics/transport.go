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
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	appDynamicsAlertEndoint := kithttp.NewServer(makeAppDynamicsAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	addAppdynamicsEndpoint := kithttp.NewServer(makeAddAppdynamicsEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	addAppdynamicsQueueEndpoint := kithttp.NewServer(makeAddAppdynamicsQueueEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	executeCommandFromAppdynamics := kithttp.NewServer(makExecuteCommandFromAppdynamics(service), decodeExecuteRequest, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/appdynamics/{chatid} appdynamics sendAlert
	//
	// Sends an appdynamics alert message to the chat id
	//
	// Use the following template definition for your appdynamics http  request
	//
	// ```
	// #macro(EntityInfo $item)
	//	"name": "${item.name}"
	// #end
	//
	//	{
	//	"environment": "DEV / SIT",
	//	"policy" : {
	//		"digestDurationInMins": "${policy.digestDurationInMins}",
	//		"name": "${policy.name}"
	//
	//	},
	//	"action": {
	//		"triggerTime":  "${action.triggerTime}",
	//		"name": "${action.name}"
	//
	//	},
	//	"events": [
	//		#foreach(${event} in ${fullEventList})
	//		{
	//			"severity": "${event.severity}",
	//			"application": {
	//				#EntityInfo($event.application)
	//			},
	//			"tier": {
	//				#EntityInfo($event.tier)
	//			},
	//			"node": {
	//				#EntityInfo($event.node)
	//			},
	//
	//			"displayName": "$event.displayName",
	//			"eventMessage": "$event.eventMessage"
	//		}
	//		#if($foreach.count != $fullEventList.size()) , #end
	//		#end
	//	]
	// }
	// ```
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
	//   description: send alert request
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/appdynamicsMessage"
	// responses:
	//   '200':
	//     description: Command has been executed successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/appdynamics/{chatid:[0-9]+}", appDynamicsAlertEndoint).Methods("POST")
	r.Handle("/api/appdynamics/{chatid:[0-9]+}/queue", addAppdynamicsEndpoint).Methods("POST")

	// swagger:operation POST /api/appdynamics/{chatid}/system appdynamics configureAppdynamics
	//
	// Configure the appdynamics endpoints for the chat, used when the bot has to look up information from appdynamics
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
	//   description: add ssh key request
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/AddAppdynamicsEndpointRequest"
	// responses:
	//   '200':
	//     description: Command has been executed successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/appdynamics/{chatid:[0-9]+}/endpoint", addAppdynamicsQueueEndpoint).Methods("POST")

	// swagger:operation POST /api/appdynamics/{chatid}/execute appdynamics executeCommand
	//
	// Executes a predefined SSH Command using the key added to your group
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
	//   description: add ssh key request
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/ExecuteAppDynamicsCommandRequest"
	// responses:
	//   '200':
	//     description: Command has been executed successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/appdynamics/{chatid:[0-9]+}/execute", executeCommandFromAppdynamics).Methods("POST")

	return r

}

func decodeExecuteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := ExecuteAppDynamicsCommandRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
