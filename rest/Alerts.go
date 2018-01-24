package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"github.com/zamedic/go2hal/service"
	"encoding/json"
	"encoding/base64"
	"errors"
)

type imageAlertMessage struct {
	Message, Image string
	InternalError bool
}


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

func imageAlertHandler(w http.ResponseWriter, r *http.Request) {
	var req imageAlertMessage
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if req.Image != "" {
		b, err := base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if req.InternalError {
			service.SendImageToHeartbeatGroup(b)
		} else {
			service.SendImageToAlertGroup(b)

		}
	}
	if req.InternalError{
		service.SendError(errors.New(req.Message))
	} else {
		service.SendAlert(req.Message)
	}
}

func sendBusinessAlert(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(err,w)
		return
	}
	err = service.SendNonTechnicalAlert(string(body))

}



func handleError(err error,w http.ResponseWriter){
	w.WriteHeader(http.StatusInternalServerError)
	service.SendError(err)
}
