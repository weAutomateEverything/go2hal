package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

/*
RouterObject provides a pointer to the underlying mux object for status checks
 */
type RouterObject struct {
	Mux *mux.Router
}

var router *RouterObject

func init() {
	router = &RouterObject{}
	go func() {
		log.Println("Starting HTTP Server...")
		log.Fatal(http.ListenAndServe(":8000", getRouter()))
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/alert", alertHandler)
	r.HandleFunc("/alert/aws/container", handleEc2ContainerAlert)
	r.HandleFunc("/github", handleGithubMessage)
	r.HandleFunc("/status", status)
	r.HandleFunc("/bot", addBot).Methods("POST")
	r.HandleFunc("/bot", botStatus).Methods("GET")
	r.HandleFunc("/httpEndpoint", addHTTPEndpoint).Methods("POST")
	r.HandleFunc("/appdynamics", receiveAppDynamicsAlert).Methods("POST")
	r.HandleFunc("/delivery", receiveDeliveryNotification).Methods("POST")
	r.HandleFunc("/recipe", addRecipe).Methods("POST")
	r.HandleFunc("/environment", addChefEnvironment).Methods("POST")
	r.HandleFunc("/chefAudit", sendAnalyticsMessage).Methods("POST")

	router.Mux = r
	return r
}

/*
Router starts the router service
 */
func Router() (*RouterObject) {
	return router
}
