package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/telegram"
	"io/ioutil"
	"log"
)

func alertHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(err,w)
		return
	}
	err = sendAlert(string(body))
	if err != nil {
		handleError(err,w)
	}
	database.ReceiveAlert()
	w.WriteHeader(http.StatusOK)
}

func sendAlert(message string) error {
	alertGroup, err := database.AlertGroup()
	if err != nil{
		return err
	}
	err = telegram.SendMessage(alertGroup, message, 0)
	return err
}

func handleError(err error,w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	log.Println(err)
}
