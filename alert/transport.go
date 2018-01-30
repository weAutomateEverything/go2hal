package alert

import (
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)


	alertHandler := kithttp.NewServer(makeAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts..., )
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service), decodeImageAlertRequest, gokit.EncodeResponse, opts..., )
	busienssAlertHandler := kithttp.NewServer(makeBusinessAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts..., )

	r := mux.NewRouter()

	r.Handle("/alert", alertHandler).Methods("POST")
	r.Handle("/alert/image", imageAlertHandler).Methods("POST")
	r.Handle("/alert/business", busienssAlertHandler).Methods("POST")

	return r
}

func decodeImageAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var a imageAlertMessage
	err := json.NewDecoder(r.Body).Decode(&a)
	return a, err
}
