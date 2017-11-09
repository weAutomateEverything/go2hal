package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"github.com/zamedic/go2hal/service"
)

func receiveAppDynamicsAlert(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		service.SendError(err)
		return
	}

	database.SaveAudit("APPDYNAMICS",string(body))
	database.ReceiveAppynamicsMessage()
	service.SendAppdynamicsAlert(string(body))

}
