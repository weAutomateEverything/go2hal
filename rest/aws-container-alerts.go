package rest

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
)

func handleEc2ContainerAlert(w http.ResponseWriter, r *http.Request) {
	var f interface{}
	body, err := ioutil.ReadAll(r.Body)
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
	subscribeUrl := m["SubscribeURL"]

	if subscribeUrl != nil {
		sendAlert("Received AWS subscription url:" + subscribeUrl.(string))
	}

	stringMsgString :=m["Message"].(string)
	var g interface{}
	err = json.Unmarshal([]byte(stringMsgString),&g)
	if err != nil {
		log.Println(err)
		return
	}
	msgObj :=  g.(map[string]interface{})
	event := msgObj["detail-type"].(string)
	if event == "ECS Task State Change"{
		detail := msgObj["detail"].(map[string]interface{})
		desired := detail["desiredStatus"].(string)
		last := detail["lastStatus"].(string)
		service := detail["group"].(string)
		sendAlert("*"+event+"*\n"+"Desired State: "+desired+"\nLast State: "+last+"\nservice: "+service)
	}
}
