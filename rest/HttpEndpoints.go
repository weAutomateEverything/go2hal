package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
	"log"
	"github.com/zamedic/go2hal/service"
)

func addHTTPEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpoint database.HTTPEndpoint
	err := json.NewDecoder(r.Body).Decode(&endpoint)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = service.CheckEndpoint(endpoint)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}
	database.AddHTMLEndpoint(endpoint)
	w.WriteHeader(http.StatusOK)

}
