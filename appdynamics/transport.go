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
	addAppdynamicsEndpoint := kithttp.NewServer(makeAddAppdynamicsEndpoint(service), decodeAddEndpointRequest, gokit.EncodeResponse, opts...)
	addAppdynamicsQueueEndpoint := kithttp.NewServer(makeAddAppdynamicsQueueEndpoint(service), decodeAddMqEndpointRequest, gokit.EncodeResponse, opts...)
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

	// swagger:operation POST /api/appdynamics/{chatid}/endpoint appdynamics configureAppdynamics
	//
	// Configure the appdynamics endpoints for the chat, used when the bot has to look up information from appdynamics
	//
	//HAL needs to know the URL, user and password to use to query appdynamics
	//
	//so, if your user is A-user, your group is customer1 and your password is secret, with app dynamics available on http://appd.yourcomany.com:8090
	//then the Endpoint address would be
	//http://A-user%40customer1:secret@appd.yourcomany.com:8090
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
	//   description: Appdynamics endpoint
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
	r.Handle("/api/appdynamics/{chatid:[0-9]+}/endpoint", addAppdynamicsEndpoint).Methods("POST")

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

func decodeAddEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := AddAppdynamicsEndpointRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}

func decodeAddMqEndpointRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := AddAppdynamicsQueueEndpointRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err
}
