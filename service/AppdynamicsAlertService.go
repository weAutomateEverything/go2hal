package service

import (
	json2 "encoding/json"
	"log"
	"strings"
	"gopkg.in/kyokomi/emoji.v1"
	"bytes"
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

		application := event["application"].(map[string]interface{})
		tier := event["tier"].(map[string]interface{})
		node := event["node"].(map[string]interface{})

		message = strings.Replace(message,"<b>","*",-1)
		message = strings.Replace(message,"</b>","*",-1)
		message = strings.Replace(message,"<br>","\n",-1)

		var buffer bytes.Buffer
		buffer.WriteString(emoji.Sprintf(":red_circle:"))
		buffer.WriteString(" ")
		buffer.WriteString(message)
		buffer.WriteString("\n")

		if application != nil {
			buffer.WriteString("*Application:* ")
			buffer.WriteString(application["name"].(string))
			buffer.WriteString("\n")
		}

		if tier != nil {
			buffer.WriteString("*Tier:* ")
			buffer.WriteString(tier["name"].(string))
			buffer.WriteString("\n")
		}

		if node != nil {
			buffer.WriteString("*Node:* ")
			buffer.WriteString(node["name"].(string))
			buffer.WriteString("\n")
		}

		log.Printf("Sending Alert %s",buffer.String())
		SendAlert(buffer.String())
	}
}
