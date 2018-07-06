package telegram

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

/*
MakeHandler returns a HTTP Restul endpoint to handle user requests
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	//requestAuth := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(),jwt.SigningMethodHS256,gokitjwt.StandardClaimsFactory)( makeTelegramAuthRequestEndpoint(service)),
	//gokit.DecodeString, gokit.EncodeResponse, opts...)
	requestAuth := kithttp.NewServer(makeTelegramAuthRequestEndpoint(service),
		decodeAuthRequest, gokit.EncodeResponse, opts...)

	authpoll := kithttp.NewServer(makeTelegramAuthPollEndpoint(service),
		decodeAuthPoll, encodeAuthResoinse)
	sendKeyboardOptions :=kithttp.NewServer(makeSendKeyboardEndpoint(service),
		decodeSendKeyBoardRequest, gokit.EncodeResponse)
	setState :=kithttp.NewServer(makeSetStateRequestEndpoint(service),
		decodeSetStateRequest, gokit.EncodeResponse)
	r := mux.NewRouter()

	r.Handle("/api/telegram/auth", requestAuth).Methods("POST")
	r.Handle("/api/telegram/auth/{id}", authpoll).Methods("GET")
	r.Handle("/api/telegram/keyboard", sendKeyboardOptions).Methods("POST")
	r.Handle("/api/telegram/state", setState).Methods("POST")
	return r

}

func decodeAuthRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	//jwt := ctx.Value(gokitjwt.JWTClaimsContextKey).(CustomClaims)
	v := authRequestObject{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func decodeAuthPoll(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	return vars["id"], nil

}
func decodeSendKeyBoardRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := sendKeyBoardRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
func decodeSetStateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := setStateRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
type errorer interface {
	error() error
}

func encodeAuthResoinse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if response == nil {
		return nil
	}
	if e, ok := response.(errorer); ok && e.error() != nil {
		gokit.EncodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Add("Authorization", "Bearer "+response.(string))
	return nil
}
