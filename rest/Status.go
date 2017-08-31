package rest

import "net/http"

func status(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
}
