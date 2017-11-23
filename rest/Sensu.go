package rest

import (
	"net/http"
	"io/ioutil"
	"github.com/zamedic/go2hal/service"
	"errors"
)

func sensuSlackAlert(w http.ResponseWriter, r *http.Request) {
	msg, _ := ioutil.ReadAll(r.Body)
	service.SendError(errors.New(string(msg)))
	w.WriteHeader(http.StatusOK)
}

