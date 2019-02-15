package ssh

import (
	"github.com/dgrijalva/jwt-go"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/weAutomateEverything/go2hal/telegram"
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

	addCommandServer := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeAddCommandEndpoint(service)), decodeAddCommand, gokit.EncodeResponse, opts...)

	addKeyServer := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeAddKeyEndpoint(service)), decodeAddKey, gokit.EncodeResponse, opts...)

	executeServer := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeExecuteCommandEndpoint(service)), decodeExeuteCommand, gokit.EncodeResponse, opts...)

	addServerserver := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeAddServerEndpoint(service)), decodeAddServer, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/ssh/key ssh addKey
	//
	// Sets the ssh private key to be used for this chat to execute ssh commands.
	// Each chat id can currently only have 1 ssh key
	//
	//
	// ---
	// security:
	// - api_key: []
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   description: add ssh key request. Ensure that the key value is base64 encoded!
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
	r.Handle("/api/ssh/key", addKeyServer).Methods("POST")

	// swagger:operation POST /api/ssh/command ssh addCommand
	//
	// Adds a command to the chat, this will allow the execute command opparation to execute the command by using
	// the name provided
	//
	//
	// ---
	// security:
	// - api_key: []
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
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
	r.Handle("/api/ssh/command", addCommandServer).Methods("POST")

	// swagger:operation POST /api/ssh/execute ssh executeCommand
	//
	// Executes a predefined SSH Command using the key added to your gropup
	//
	// ---
	// security:
	// - api_key: []
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
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
	r.Handle("/api/ssh/execute/", executeServer).Methods("POST")

	// swagger:operation POST /api/ssh/server ssh addServer
	//
	// Adds a new server that the user can select when executing a command
	//
	//
	// ---
	// security:
	// - api_key: []
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   description: add server request
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/addServer"
	// responses:
	//   '200':
	//     description: Command has been added
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/ssh/server", addServerserver).Methods("POST")

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

func decodeAddServer(_ context.Context, r *http.Request) (interface{}, error) {
	var execute addServer
	err := json.NewDecoder(r.Body).Decode(&execute)
	return execute, err
}
