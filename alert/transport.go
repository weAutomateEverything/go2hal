package alert

import (
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"encoding/json"
	"context"
)

/*
MakeHandler returns a rest http handler to send alerts.

The machine learning service can be set to nil if you do not wish to save the requests.
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	alertHandler := kithttp.NewServer(makeAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	recipeKeyboardHandler := kithttp.NewServer(makeKeyboardRecipeAlertEndpoint(service), decodeKeyboardAlertRequest, gokit.EncodeResponse, opts...)
	environmentKeyboardHandler := kithttp.NewServer(makeEnvironmentAlertEndpoint(service), decodeKeyboardAlertRequest, gokit.EncodeResponse, opts...)
	nodesKeyboardHandler := kithttp.NewServer(makeNodesAlertEndpoint(service), decodeKeyboardAlertRequest, gokit.EncodeResponse, opts...)
	imageAlertHandler := kithttp.NewServer(makeImageAlertEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)
	heartbeatAlertHandler := kithttp.NewServer(makeHeartbeatMessageEncpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	heartbeatImageHandler := kithttp.NewServer(makeImageHeartbeatEndpoint(service), gokit.DecodeFromBase64, gokit.EncodeResponse, opts...)
	busienssAlertHandler := kithttp.NewServer(makeBusinessAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)
	alertErrorHandler := kithttp.NewServer(makeAlertErrorHandler(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	r := mux.NewRouter()

	r.Handle("/alert/", alertHandler).Methods("POST")
	r.Handle("/alert/image", imageAlertHandler).Methods("POST")
	r.Handle("/alert/keyboard/recipe", recipeKeyboardHandler).Methods("POST")
	r.Handle("/alert/keyboard/environment", environmentKeyboardHandler).Methods("POST")
	r.Handle("/alert/keyboard/nodes", nodesKeyboardHandler).Methods("POST")
	r.Handle("/alert/heartbeat", heartbeatAlertHandler).Methods("POST")
	r.Handle("/alert/heartbeat/image", heartbeatImageHandler).Methods("POST")

	r.Handle("/alert/error", alertErrorHandler).Methods("POST")

	r.Handle("/alert/business", busienssAlertHandler).Methods("POST")

	return r
}
func decodeKeyboardAlertRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request KeyboardAlertRequest
	if err := json.NewDecoder(r.Body).Decode(&request.Nodes); err != nil {
		return nil, err
	}
	return request, nil

}
