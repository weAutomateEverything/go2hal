package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"github.com/zamedic/go2hal/database"
)

func receiveDeliveryNotification(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	database.SaveAudit("APPDYNAMICS",string(body))
	database.ReceiveChefDeliveryMessage()
}