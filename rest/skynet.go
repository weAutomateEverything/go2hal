package rest

import (
	"net/http"
	"io/ioutil"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/service"
	"encoding/json"
	"log"
	"strings"
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
	s, _ := ioutil.ReadAll(r.Body)
	log.Printf("Received rebuild request %s",s)
	err := json.NewDecoder(strings.NewReader(string(s))).Decode(&rebuild)
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
