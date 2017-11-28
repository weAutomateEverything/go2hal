package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"github.com/zamedic/go2hal/service"
	"encoding/json"
	"log"
	"fmt"
)

type appdynamicsQueueEndpoint struct {
	Name        string
	Application string
	Metricpath  string
}

type appdynamicsEndpoint struct {
	Endpoint string
}

type executeAppDynamicsCommand struct {
	CommandName, NodeID, ApplicationID string
}


type businessAlert struct {
	Severity, Type, DisplayName, SummaryMessage string
}

func receiveAppDynamicsAlert(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		service.SendError(err)
		return
	}

	database.SaveAudit("APPDYNAMICS", string(body))
	database.IncreaseValue("APPDYNAMICS_ALERT_REQUESTS")
	service.SendAppdynamicsAlert(string(body))
}

func addAppDynamicsEndpoint(w http.ResponseWriter, r *http.Request) {
	e := appdynamicsEndpoint{}
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = service.AddAppdynamicsEndpoint(e.Endpoint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)

}

func addAppdynamicsQueueEndpoint(w http.ResponseWriter, r *http.Request) {
	var appDynamics appdynamicsQueueEndpoint
	err := json.NewDecoder(r.Body).Decode(&appDynamics)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = service.AddAppDynamicsQueue(appDynamics.Name, appDynamics.Application, appDynamics.Metricpath)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func executeCommandFromAppdynamics(w http.ResponseWriter, r *http.Request) {
	c := executeAppDynamicsCommand{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		s, _ := ioutil.ReadAll(r.Body)
		service.SendError(fmt.Errorf("received a bad request to execute a command. %s", s))
		service.SendError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = service.ExecuteCommandFromAppd(c.CommandName,c.ApplicationID,c.NodeID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	database.IncreaseValue("APPDYNAMICS_EXECUTE_REQUESTS")
	w.WriteHeader(http.StatusOK)
}

func businessAppDynamicsAlert(w http.ResponseWriter, r *http.Request) {
	c := businessAlert{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		s, _ := ioutil.ReadAll(r.Body)
		service.SendError(fmt.Errorf("received a bad request to busines alert service. %s", s))
		service.SendError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	service.SendError(fmt.Errorf("business alert. Servity: %s - type: %s - Display Name: %s - Summary Message %s", c.Severity, c.Type, c.DisplayName, c.SummaryMessage))
	database.IncreaseValue("APPDYNAMICS_BUSINESS_ALERT_REQUESTS")
	w.WriteHeader(http.StatusOK)
}
