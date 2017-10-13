package service

import (
	json2 "encoding/json"
	"log"
	"strings"
)

/*
SendAppdynamicsAlert will send an alert message from the data contained in the even message field
 */
func SendAppdynamicsAlert(message string) {
	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", message)
		return
	}

	events := dat["events"].([]interface{})
	for _, event := range events{
		event := event.(map[string]interface{})

		message := event["eventMessage"].(string)
		message = strings.Replace(message,"<b>","*",-1)
		message = strings.Replace(message,"</b>","*",-1)

		message = strings.Replace(message,"<br>","\n",-1)

		log.Printf("Sending Alert %s",message)
		SendAlert(message)
	}
}
