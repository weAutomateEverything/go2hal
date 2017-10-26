package rest

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/zamedic/go2hal/service"
)

type appdynamicsQueueEndpoint struct {
	Name       string
	Endpoint   string
	Metricpath string
}

func addAppdynamicsQueueEndpoint(w http.ResponseWriter, r *http.Request) {
	var appDynamics appdynamicsQueueEndpoint
	err := json.NewDecoder(r.Body).Decode(&appDynamics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = service.AddAppDynamicsQueue(appDynamics.Name, appDynamics.Endpoint, appDynamics.Metricpath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}
