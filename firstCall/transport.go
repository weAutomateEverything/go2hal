package firstCall

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"github.com/weAutomateEverything/go2hal/telegram"
	"io/ioutil"
	"net/http"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service CalloutFunction, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	getEndpoints := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeGetDefaultCalloutEndpoint(service)), gokit.DecodeString, gokit.EncodeResponse,
		opts...)

	setEndpoints := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeSetDefaultCalloutEndpoint(service.(DefaultCalloutService))), decodeAddEndpointRequest, gokit.EncodeResponse,
		opts...)

	r := mux.NewRouter()

	r.Handle("/api/firstcall/defaultCallout", setEndpoints).Methods("POST")
	r.Handle("/api/firstcall/defaultCallout", getEndpoints).Methods("GET")

	return r

}

func decodeAddEndpointRequest(ctx context.Context, r *http.Request) (resp interface{}, err error) {
	req := defaultCalloutRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &req)
	resp = req
	return

}
