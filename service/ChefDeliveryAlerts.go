package service

import (
	json2 "encoding/json"
	"strings"
	"log"
	"bytes"
)

/*
SendDeliveryAlert will unmarshal the input json and send an alert to the telegram group.
 */
func SendDeliveryAlert(message string) {
	var dat map[string]interface{}

	message = strings.Replace(message, "\n", "\\n", -1)

	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", message)
	}

	attachments := dat["attachments"].([]interface{})
	//Loop though the attachmanets, there should be only 1
	var buffer bytes.Buffer
	buffer.WriteString("*Chef Delivery*\n")
	for _, attachment := range attachments {
		attachment_i := attachment.(map[string]interface{})
		fields := attachment_i["fields"].([]interface{})

		//Loop through the fields
		for _, field := range fields {
			field_i := field.(map[string]interface{})
			buffer.WriteString("*")
			buffer.WriteString(field_i["title"].(string))
			buffer.WriteString("*")
			buffer.WriteString(field_i["value"].(string))
			buffer.WriteString("\n")
		}
	}
	log.Printf("Sending Alert: %s",buffer.String())

	body := dat["text"].(string)
	bodies := strings.Split(body,"\n");

	url := bodies[0]
	url = strings.Replace(url,"<","",-1)
	url = strings.Replace(url,">","",-1)

	parts := strings.Split(url,"|")

	buffer.WriteString("[")
	buffer.WriteString(parts[1])
	buffer.WriteString(" - ")
	buffer.WriteString(bodies[1])
	buffer.WriteString("](")
	buffer.WriteString(parts[0])
	buffer.WriteString(")")

	SendAlert(buffer.String())

}
