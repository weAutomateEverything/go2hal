package service

import (
	json2 "encoding/json"
	"log"
	"bytes"
	"github.com/zamedic/go2hal/database"
	"strings"
)
/*
SendAnalyticsAlert will check if we have a chef recipe configured for the alert. If we do, it will send an alert.
 */
func SendAnalyticsAlert(message string) {
	if(!checkSend(message)){
		log.Println("Ignoreing message: "+message)
		return
	}
	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", message)
		return
	}

	attachments := dat["attachments"].([]interface{})
	//Loop though the attachmanets, there should be only 1

	var buffer bytes.Buffer
	buffer.WriteString("*Analytics Event*\n")
	buffer.WriteString(dat["text"].(string))
	buffer.WriteString("\n")
	getfield(attachments, &buffer)

	log.Printf("Sending Alert: %s", buffer.String())
	SendAlert(buffer.String())
}

func checkSend(message string) bool {
	message = strings.ToUpper(message)
	log.Printf("Checking if we should send: %s",message)
	recipes := database.GetRecipes()
	for _, recipe := range recipes {
		check := "RECIPE["+strings.ToUpper(recipe)+"]"
		log.Printf("Comparing %s",check)
		if strings.Contains(message,check) {
			log.Printf("Match Found, returning true")
			return true
		}
	}
	log.Printf("No match found, not sending message")
	return false;
}
