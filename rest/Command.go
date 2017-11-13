package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
	"fmt"
	"io/ioutil"
)

type command struct {
	Name, Command string
}

type executeAppDynamicsCommand struct {
	CommandName, NodeID string
}

func addCommand(w http.ResponseWriter, r *http.Request) {
	c := command{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = database.AddCommand(c.Name, c.Command)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

}
