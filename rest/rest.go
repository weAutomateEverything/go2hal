package rest

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"errors"
	"github.com/zamedic/go2hal/service"
	"runtime/debug"
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			service.SendError(errors.New(fmt.Sprint(err)))
			service.SendError(errors.New(string(debug.Stack())))


		}
	}()
}

func getRouter() *mux.Router {


	r.HandleFunc("/github", handleGithubMessage)
	r.HandleFunc("/status", status)
	r.HandleFunc("/bot", addBot).Methods("POST")
	r.HandleFunc("/bot", botStatus).Methods("GET")
	r.HandleFunc("/httpEndpoint", addHTTPEndpoint).Methods("POST")

	r.HandleFunc("/appdynamics", receiveAppDynamicsAlert).Methods("POST")
	r.HandleFunc("/appdynamics/system/queue", addAppdynamicsQueueEndpoint).Methods("POST")
	r.HandleFunc("/appdynamics/system", addAppDynamicsEndpoint).Methods("POST")
	r.HandleFunc("/appdynamics/execute", executeCommandFromAppdynamics).Methods("POST")
	r.HandleFunc("/appdynamics/alert/nontech", businessAppDynamicsAlert).Methods("POST")

	r.HandleFunc("/delivery", receiveDeliveryNotification).Methods("POST")

	r.HandleFunc("/chefAudit", sendAnalyticsMessage).Methods("POST")

	r.HandleFunc("/skynet", sendSkynetAlert).Methods("POST")
	r.HandleFunc("/skynet/endpoint", addSkynetEndpoint).Methods("POST")
	r.HandleFunc("/skynet/rebuild",rebuildNode).Methods("POST")

	r.HandleFunc("/command",addCommand).Methods("POST")
	r.HandleFunc("/command/key",addKey).Methods("POST")

	r.HandleFunc("/selenium",addSeleniumCheck).Methods("POST")

	r.HandleFunc("/sensu", sensuSlackAlert).Methods("POST")

	r.HandleFunc("/config/jira",saveJira).Methods("POST")
	r.HandleFunc("/config/callout",saveCallout).Methods("POST")
	r.HandleFunc("/config/chef",addChefClient).Methods("POST")
	r.HandleFunc("/config/chef/recipe",addRecipe).Methods("POST")
	r.HandleFunc("/config/chef/environment",addChefEnvironment).Methods("POST")


	r.HandleFunc("/callout",invokeCallout).Methods("POST")

	r.HandleFunc("/users",addUser).Methods("POST")
	
	router.Mux = r
	return r
}

/*
Router starts the router service
 */
func Router() (*RouterObject) {
	return router
}
