package analytics

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"net/http"
)

//MakeHandler returns a Rest HTTP endpoint for chef alerts
//Machine learning can be left as nil if you do not wish to store the request
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	analyticsHandler := kithttp.NewServer(makeAnalyticsEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/api/chefAudit", analyticsHandler).Methods("POST")

	return r

}
