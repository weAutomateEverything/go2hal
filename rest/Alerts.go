package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"log"
	"github.com/zamedic/go2hal/service"
)

func alertHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(err,w)
		return
	}
	err = service.SendAlert(string(body))
	if err != nil {
		handleError(err,w)
	}
	database.ReceiveAlert()
	w.WriteHeader(http.StatusOK)
}



func handleError(err error,w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
