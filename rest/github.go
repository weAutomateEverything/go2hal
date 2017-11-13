package rest

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
	"github.com/zamedic/go2hal/service"
	"github.com/zamedic/go2hal/database"
)

func handleGithubMessage(w http.ResponseWriter, r *http.Request) {
	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
	database.SaveAudit("GITHUB",string(body))
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Println(err)
		return
	}
	m := f.(map[string]interface{})
	status := m["state"].(string)
	description := m["description"].(string)

	result := fmt.Sprintf("*GITHUB*\n %s - %s",status,description)
	service.SendAlert(result)
	w.WriteHeader(http.StatusOK)

}

