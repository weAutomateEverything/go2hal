package chef

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
	"net/http"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := gokit.GetServerOpts(logger)

	chefDeliveryEndpoint := kithttp.NewServer(makeChefDeliveryAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/delivery", chefDeliveryEndpoint).Methods("POST")

	return r

}
