package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
)

func sendAnalyticsMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	database.ReceiveChefAuditMessage()
	service.SendAnalyticsAlert(string(body))
}
