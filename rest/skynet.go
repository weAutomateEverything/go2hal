package rest

import (
	"net/http"
	"io/ioutil"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
	"encoding/json"
)

type skynet struct {
	URL, Username, Password string
}

type skynetRebuild struct {
	NodeName string `json:"Nodename"`
	User     string `json:"User"`
}

func sendSkynetAlert(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		service.SendError(err)
	}
	database.ReceiveSkynetMessage()
	service.SendSkynetAlert(string(body))
}

func addSkynetEndpoint(w http.ResponseWriter, r *http.Request) {
	var skynet skynet
	err := json.NewDecoder(r.Body).Decode(&skynet)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	database.AddSkynetEndpoint(skynet.URL, skynet.Username, skynet.Password)
	w.WriteHeader(http.StatusOK)
}

func rebuildNode(w http.ResponseWriter, r *http.Request) {
	var rebuild skynetRebuild
	err := json.NewDecoder(r.Body).Decode(&rebuild)
	if err != nil {
		service.SendError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	go func() {
		service.RecreateNode(rebuild.NodeName, rebuild.User)
	}()
}
