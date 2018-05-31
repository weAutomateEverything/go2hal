package alert

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"

	"context"
	"encoding/base64"
	"io/ioutil"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	alertHandler := kithttp.NewServer(makeAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)
	documentAlertHandler := kithttp.NewServer(makeDocumentAlertEndpoint(service), decodeSendDocumentRequest, gokit.EncodeResponse, opts...)
	alertErrorHandler := kithttp.NewServer(makeAlertErrorHandler(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api//alert/{chatid} alert sendTextAlert
	//
	// Send a text alert to a telegram group
	//
	//
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// - name: message
	//   in: body
	//   description: message you want to send
	//   required: true
	//   schema:
	//     type: string
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}", alertHandler).Methods("POST")

	// swagger:operation POST /api/alert/{chatid}/image alert sendImageAlert
	//
	// Send a image to a telegram group
	//
	//
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// - name: message
	//   in: body
	//   description: message you want to send enocded in base64 format
	//   required: true
	//   schema:
	//     type: string
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}/image", imageAlertHandler).Methods("POST")

	// swagger:operation POST /api/alert/{chatid}/document/{extension} alert sendDocumentAlert
	//
	// Send a text alert to a telegram group
	//
	//
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// - name: extension
	//   in: path
	//   description: document extension (pdf, doc, xls)
	//   required: true
	//   type: string
	// - name: message
	//   in: body
	//   description: raw document data encoded in base64
	//   required: true
	//   schema:
	//     type: string
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}/document/{extension}", documentAlertHandler).Methods("POST")

	r.Handle("/api/alert/error", alertErrorHandler).Methods("POST")

	return r
}

func decodeSendDocumentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	base64msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	_, err = base64.StdEncoding.Decode(base64msg, base64msg)
	if err != nil {
		return nil, err
	}

	return sendDocumentRequest{
		exension: vars["extension"],
		document: base64msg,
	}, nil

}
