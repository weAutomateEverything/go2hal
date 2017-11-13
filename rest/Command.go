package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
)

type command struct {
	Name, Command string
}

type key struct {
	Username, Key string
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

func addKey(w http.ResponseWriter, r *http.Request) {
	c := key{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = database.AddKey(c.Username, c.Key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}


