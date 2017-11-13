package rest

import (
	"net/http"
	"io/ioutil"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
)

func sendAnalyticsMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		service.SendError(err)
		return
	}
	database.ReceiveChefAuditMessage()
	service.SendAnalyticsAlert(string(body))
}
