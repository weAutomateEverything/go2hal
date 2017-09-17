package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/service"
)

type chefClient struct{
	Name, Key, URL string
}

func addChefClient(w http.ResponseWriter, r *http.Request) {
	var chef chefClient
	err := json.NewDecoder(r.Body).Decode(&chef)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = service.AddChefClient(chef.Name,chef.Key,chef.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
