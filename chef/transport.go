package chef

import (
	gokitjwt "github.com/go-kit/kit/auth/jwt"
	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"

	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"github.com/weAutomateEverything/go2hal/telegram"
	"net/http"
)

//MakeHandler returns a restful http handler for the chef delivery service
//the Machine Learning service can be set to nil if you do not wish to log the http requests

func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	chefDeliveryEndpoint := kithttp.NewServer(makeChefDeliveryAlertEndpoint(service), gokit.DecodeString, gokit.EncodeResponse, opts...)

	addChefRecipeToGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeAddRecipeToGroupEndpoint(service)), decodeAddChefRequest, gokit.EncodeResponse, opts...)

	getChefRecipesForGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeGetAllGrouRecipesEndpoint(service)), gokit.DecodeString, gokit.EncodeResponse, opts...)

	addEnvironmentToGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeAddEnvironmentToGroupEndpoint(service)), decodeAddEnvironmentfRequest, gokit.EncodeResponse, opts...)

	getEnvironmentForGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(), jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)(makeGetEnvironmentForGroupEndpoint(service)), gokit.DecodeString, gokit.EncodeResponse, opts...)
	getChefRecipesByGroup := kithttp.NewServer(makeGetChefRecipesByGroupEndpoint(service), handleRequest, gokit.EncodeResponse, opts...)
	getChefEnvironmentsByGroup := kithttp.NewServer(makeGetChefEnvironmentsByGroupEndpoint(service), handleRequest, gokit.EncodeResponse, opts...)
	getChefNodes := kithttp.NewServer(makeGetChefNodesEndpoint(service), decodeGetNodesRequest, gokit.EncodeResponse, opts...)
	r := mux.NewRouter()

	r.Handle("/api/chef/delivery/{chatid:[0-9]+}", chefDeliveryEndpoint).Methods("POST")
	r.Handle("/api/chef/recipe", addChefRecipeToGroup).Methods("POST")
	r.Handle("/api/chef/recipes", getChefRecipesForGroup).Methods("GET")
	r.Handle("/api/chef/environment", addEnvironmentToGroup).Methods("POST")
	r.Handle("/api/chef/environments", getEnvironmentForGroup).Methods("GET")
	r.Handle("/api/chef/environments/group/{groupid}", getChefEnvironmentsByGroup).Methods("GET")
	r.Handle("/api/chef/nodes", getChefNodes).Methods("POST")
	r.Handle("/api/chef/recipes/group/{groupid}", getChefRecipesByGroup).Methods("GET")
	return r

}

func decodeAddChefRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var q = &addRecipeRequest{}
	err := json.NewDecoder(r.Body).Decode(&q)
	return q, err
}

func decodeAddEnvironmentfRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var q = &addEnvironmentRequest{}
	err := json.NewDecoder(r.Body).Decode(&q)
	return q, err
}
func handleRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var q = mux.Vars(r)
	return q["groupid"],nil
}
func decodeGetNodesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var q = &chefNodeRequest{}
	err := json.NewDecoder(r.Body).Decode(&q)
	return q, err
}