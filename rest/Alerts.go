package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
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
	database.IncreaseValue("ALERT_REQUESTS")
	w.WriteHeader(http.StatusOK)
}

func sendBusinessAlert(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(err,w)
		return
	}
	err = service.SendNonTechnicalAlert(string(body))
	database.IncreaseValue("BUSINESS_ALERT_REQUESTS")

}



func handleError(err error,w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	service.SendError(err)
}
