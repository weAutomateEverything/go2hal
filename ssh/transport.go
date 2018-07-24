package ssh

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
)

//MakeHandler retuns a http rest request handler for sensu
//the machine learning service can be nil if you do not wish to save the request message
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	addCommand := makeAddCommandEndpoint(service)
	addKey := makeAddKeyEndpoint(service)
	execute := makeExecuteCommandEndpoint(service)

	addCommandServer := kithttp.NewServer(addCommand, decodeAddCommand, gokit.EncodeResponse, opts...)
	addKeyServer := kithttp.NewServer(addKey, decodeAddKey, gokit.EncodeResponse, opts...)
	executeServer := kithttp.NewServer(execute, decodeExeuteCommand, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/ssh/key/{chatid} ssh addKey
	//
	// Sets the ssh private key to be used for this chat to execute ssh commands.
	// Each chat id can currently only have 1 ssh key
	//
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
	//     "$ref": "#/definitions/addKey"
	// responses:
	//   '200':
	//     description: Key has been added
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/ssh/key/{chatid:[0-9]+}", addKeyServer).Methods("POST")

	// swagger:operation POST /api/ssh/command/{chatid} ssh addCommand
	//
	// Adds a command to the chat, this will allow the execute command opparation to execute the command by using
	// the name provided
	//
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
	//     "$ref": "#/definitions/addCommand"
	// responses:
	//   '200':
	//     description: Command has been added
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/ssh/command/{chatid:[0-9]+}", addCommandServer).Methods("POST")

	// swagger:operation POST /api/ssh/execute/{chatid} ssh executeCommand
	//
	// Executes a predefined SSH Command using the key added to your gropup
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
	//     "$ref": "#/definitions/executeCommand"
	// responses:
	//   '200':
	//     description: Command has been executed successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/ssh/execute/{chatid:[0-9]+}", executeServer).Methods("POST")

	return r
}

func decodeAddKey(_ context.Context, r *http.Request) (interface{}, error) {
	var addKey addKey
	err := json.NewDecoder(r.Body).Decode(&addKey)
	return addKey, err
}

func decodeAddCommand(_ context.Context, r *http.Request) (interface{}, error) {
	var addCommand addCommand
	err := json.NewDecoder(r.Body).Decode(&addCommand)
	return addCommand, err
}

func decodeExeuteCommand(_ context.Context, r *http.Request) (interface{}, error) {
	var execute executeCommand
	err := json.NewDecoder(r.Body).Decode(&execute)
	return execute, err
}
