package service

import (
	json2 "encoding/json"
	"log"
)

func SendAppdynamicsAlert(message string) {
	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", message)
		return
	}

	events := dat["events"].([]interface{})
	for _, event := range events{
		event := event.(map[string]interface{})
		log.Printf("Sending Alert %s",event["eventMessage"].(string))
		SendAlert(event["eventMessage"].(string))
	}
}
