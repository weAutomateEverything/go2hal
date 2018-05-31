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

	r.Handle("/api/httpEndpoints", addEndpoint).Methods("POST")
	r.Handle("/api/httpEndpoints", getEndpoints).Methods("GET")
	//r.Handle("/httpendpoints/{id}", authpoll).Methods("DELETE")

	return r

}

func decodeAddEndpointRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := addHttpRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
