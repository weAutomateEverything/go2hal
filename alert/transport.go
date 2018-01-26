package alert

import (
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"context"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	alertHandler := kithttp.NewServer(makeAlertEndpoint(service),decodeAlertRequest,encodeResponse, opts...,)
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service),decodeImageAlertRequest,encodeResponse, opts...,)
	busienssAlertHandler := kithttp.NewServer(makeBusinessAlertEndpoint(service),decodeAlertRequest,encodeResponse, opts...,)


	r := mux.NewRouter()

	r.Handle("/alert", alertHandler).Methods("POST")
	r.Handle("/alert/image", imageAlertHandler).Methods("POST")
	r.Handle("/alert/business", busienssAlertHandler).Methods("POST")


	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

}

func decodeAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return ioutil.ReadAll(r.Body)
}

func decodeImageAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var a imageAlertMessage
	err := json.NewDecoder(r.Body).Decode(&a)
	return a, err
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if response == nil {
		return nil
	}
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}
