package rest

import (
	"net/http"
	"io/ioutil"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
)

func receiveDeliveryNotification(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		service.SendError(err)
		return
	}
	database.IncreaseValue("CHEF_DELIVERY_REQUESTS")
	service.SendDeliveryAlert(string(body))
}