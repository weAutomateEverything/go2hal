package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
	"fmt"
	"log"
)

type httpEndpoint struct {
	Name         string
	HTTPEndpoint string
}

func addHTTPEndpoint(w http.ResponseWriter, r *http.Request) {
	var endpoint httpEndpoint
	err := json.NewDecoder(r.Body).Decode(&endpoint)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	resp, err := http.Get(endpoint.HTTPEndpoint)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	if resp.StatusCode == 200 {
		database.AddHTMLEndpoint(endpoint.Name, endpoint.HTTPEndpoint)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(fmt.Sprintf("invalid response code received %d", resp.StatusCode)))
	}
}
