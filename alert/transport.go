package alert

import (
	"encoding/json"
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
	alertWithReply := kithttp.NewServer(makeReplyAlertEndpoint(service), decodeSendAlertWithReplyRequest, gokit.EncodeResponse, opts...)
	alertGetReply := kithttp.NewServer(makeGetRepliesEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	acknowledgeReply := kithttp.NewServer(makeAcknowledgeReplyEndpoint(service), decodeDeleteReplyRequest, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	// swagger:operation POST /api/alert/{chatid} alert sendTextAlert
	//
	// Send a text alert to a telegram group.
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
	// Send a image to a telegram group.
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
	// Send a text alert to a telegram group.
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

	// swagger:operation POST /api/alert/error alert sendError
	//
	// Send an error to the error group.
	//
	//
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - application/json
	// parameters:
	// - name: message
	//   in: body
	//   description: error message
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
	r.Handle("/api/alert/error", alertErrorHandler).Methods("POST")

	// swagger:operation POST /api/alert/{chatid}/withreply alert SendWithReply
	//
	// Send a text alert to a telegram group. Should anyone reply, the reply will be stored and can be retrieved.
	//
	//
	// ---
	// consumes:
	// - text/jsom
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
	//   description: request
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/sendReplyAlertMessageRequest"
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//     schema:
	//       "$ref": "#/definitions/SendReplyAlertMessageResponse"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}/withreply", alertWithReply).Methods("POST")

	// swagger:operation GET /api/alert/{chatid}/replies alert GetReplies
	//
	// Returns all the unacknowledged replies for a chat.
	//
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// responses:
	//   '200':
	//     description: Message Sent successfully
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/Replies"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}/replies", alertGetReply).Methods("GET")

	// swagger:operation GET /api/alert/{chatid}/reply/{id} alert acknowledgeReply
	//
	// Acknoeledges a reply. This reply will no longer be sent with the get replies operation.
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: chatid
	//   in: path
	//   description: chat id
	//   required: true
	//   type: integer
	// - name: id
	//   in: path
	//   description: reply ID
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: Message deleted successfully
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/errorResponse"
	r.Handle("/api/alert/{chatid:[0-9]+}/reply/{id}", acknowledgeReply).Methods("DELETE")

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

func decodeDeleteReplyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	request := acknowledgeReplyRequest{vars["id"]}
	return request, nil
}

func decodeSendAlertWithReplyRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SendReplyAlertMessageRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	return request, err

}
