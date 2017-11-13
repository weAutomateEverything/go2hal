package rest

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/zamedic/go2hal/database"
)

type environment struct {
	Environment string
}

func addChefEnvironment(w http.ResponseWriter, r *http.Request) {
	var environment environment
	err := json.NewDecoder(r.Body).Decode(&environment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.AddChefEnvironment(environment.Environment)
	w.WriteHeader(http.StatusOK)
}
