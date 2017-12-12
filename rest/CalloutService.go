package rest

import (
	"net/http"
	"encoding/json"

	"github.com/zamedic/go2hal/service"
)

type jirareq struct {
	Title, Description string
}

func invokeCallout(w http.ResponseWriter, r *http.Request){

	var req jirareq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	service.InvokeCallout(req.Title,req.Description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}