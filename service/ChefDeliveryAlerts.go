package service

import (
	json2 "encoding/json"
	"strings"
	"log"
	"bytes"
	"gopkg.in/kyokomi/emoji.v1"
	"fmt"
)

/*
SendDeliveryAlert will unmarshal the input json and send an alert to the telegram group.
 */
func SendDeliveryAlert(message string) {
	var dat map[string]interface{}

	message = strings.Replace(message, "\n", "\\n", -1)

	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		SendError(fmt.Errorf("delivery - error unmarshalling: %s", message))
		return
	}

	attachments := dat["attachments"].([]interface{})

	body := dat["text"].(string)
	bodies := strings.Split(body, "\n");
	url := bodies[0]
	url = strings.Replace(url, "<", "", -1)
	url = strings.Replace(url, ">", "", -1)

	parts := strings.Split(url, "|")

	//Loop though the attachmanets, there should be only 1
	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprint(":truck:"))
	buffer.WriteString(" ")
	buffer.WriteString("*Chef Delivery*\n")

	if (strings.Contains(bodies[1], "failed")) {
		buffer.WriteString(emoji.Sprint(":interrobang:"))

	} else {
		switch bodies[1] {
		case "Delivered stage has completed for this change.":
			buffer.WriteString(emoji.Sprint(":+1:"))

		case "Change Delivered!":
			buffer.WriteString(emoji.Sprint(":white_check_mark:"))

		case "Acceptance Passed. Change is ready for delivery.":
			buffer.WriteString(emoji.Sprint(":ok_hand:"))

		case "Change Approved!":
			buffer.WriteString(emoji.Sprint(":white_check_mark:"))

		case "Verify Passed. Change is ready for review.":
			buffer.WriteString(emoji.Sprint(":mag_right:"))
		}
	}
	buffer.WriteString(" ")

	buffer.WriteString(bodies[1])
	buffer.WriteString("\n")

	getfield(attachments, &buffer)

	buffer.WriteString("[")
	buffer.WriteString(parts[1])

	buffer.WriteString("](")
	buffer.WriteString(parts[0])
	buffer.WriteString(")")

	log.Printf("Sending Alert: %s", buffer.String())

	SendAlert(buffer.String())

}
func getfield(attachments []interface{}, buffer *bytes.Buffer) {
	for _, attachment := range attachments {
		attachmentI := attachment.(map[string]interface{})
		fields := attachmentI["fields"].([]interface{})

		//Loop through the fields
		for _, field := range fields {
			fieldI := field.(map[string]interface{})
			buffer.WriteString("*")
			buffer.WriteString(fieldI["title"].(string))
			buffer.WriteString("* ")
			buffer.WriteString(fieldI["value"].(string))
			buffer.WriteString("\n")
		}
	}
}
