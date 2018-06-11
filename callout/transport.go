package callout

import (
	"context"
	"encoding/json"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"io/ioutil"
	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	calloutHandler := kithttp.NewServer(makeCalloutEndpoint(service), decodeCalloutRequest, gokit.EncodeResponse, opts...)
	r := mux.NewRouter()

	// swagger:operation POST /api/callout/{chatid} invokeCallout
	//
	// Invokes callout by sending a telegram message to the telegram group specified by the chat id.
	// If JIRA has been configured, a JIRA ticket will be created
	// If CALLOUT has been defined, then the bot will invoke callout via alexa
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
	//   description: message you want to send enocded in base64 format
	//   required: true
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/SendCalloutRequest"
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/callout/{chatid:[0-9]+}", calloutHandler).Methods("POST")

	return r
}

func decodeCalloutRequest(_ context.Context, r *http.Request) (resp interface{}, err error) {
	res := &SendCalloutRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
