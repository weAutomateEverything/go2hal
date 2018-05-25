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

	r.Handle("/alert/{chatid:[0-9]+}", alertHandler).Methods("POST")
	r.Handle("/alert/{chatid:[0-9]+}/image", imageAlertHandler).Methods("POST")
	r.Handle("/alert/{chatid:[0-9]+}/document/{extension}", documentAlertHandler).Methods("POST")

	r.Handle("/alert/error", alertErrorHandler).Methods("POST")

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
