package chef
import (
	"net/http"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/zamedic/go2hal/gokit"
)

func MakeHandler(service Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(gokit.EncodeError),
	}

	chefDeliveryEndpoint := kithttp.NewServer(makeChefDeliveryAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/delivery", chefDeliveryEndpoint).Methods("POST")


	return r

}