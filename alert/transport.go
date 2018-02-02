package alert

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)

	alertHandler := kithttp.NewServer(makeAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)

	heartbeatAlertHandler := kithttp.NewServer(makeHeartbeatMessageEncpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	heartbeatImageHandler := kithttp.NewServer(makeImageHeartbeatEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)

	busienssAlertHandler := kithttp.NewServer(makeBusinessAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	alertErrorHandler := kithttp.NewServer(makeAlertErrorHandler(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/alert/", alertHandler).Methods("POST")
	r.Handle("/alert/image", imageAlertHandler).Methods("POST")

	r.Handle("/alert/heartbeat", heartbeatAlertHandler).Methods("POST")
	r.Handle("/alert/heartbeat/image", heartbeatImageHandler).Methods("POST")

	r.Handle("/alert/error", alertErrorHandler).Methods("POST")

	r.Handle("/alert/business", busienssAlertHandler).Methods("POST")

	return r
}
