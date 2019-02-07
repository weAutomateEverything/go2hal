package httpSmoke

import (
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"github.com/weAutomateEverything/go2hal/telegram"
	"net/http"
)

/*
MakeHandler returns a HTTP Restul endpoint to handle user requests
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	getEndpoints := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(getHTTPForGroupEndpoint(service)), gokit.DecodeString, gokit.EncodeResponse,
		opts...)

	addEndpoint := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(addHTTPEndpoint(service)), decodeAddEndpointRequest, gokit.EncodeResponse,
		opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/httpEndpoints http-checks AddEndpoint
	//
	// Adds a new HTTP Endpoint to monitor
	//
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// security:
	// - api_key: []
	// parameters:
	// - name: message
	//   in: body
	//   description: HTTP Endpoint
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/AddHttpRequest"
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/httpEndpoints", addEndpoint).Methods("POST")

	// swagger:operation GET /api/httpEndpoints http-checks GetEndpoints
	//
	// get HTTP endpoints
	//
	// ---
	// produces:
	// - application/json
	// security:
	// - api_key: []
	// responses:
	//   '200':
	//     description: success
	//     type: array
	//     items:
	//        "$ref": "#/definitions/AddHttpRequest"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/httpEndpoints", getEndpoints).Methods("GET")
	//r.Handle("/httpendpoints/{id}", authpoll).Methods("DELETE")

	return r

}

func decodeAddEndpointRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := AddHttpRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
