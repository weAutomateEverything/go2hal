package analytics

import (
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)

	analyticsHandler := kithttp.NewServer(makeAnalyticsEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/chefAudit", analyticsHandler).Methods("POST")

	return r

}
