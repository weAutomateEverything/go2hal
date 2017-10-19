package service

import (
	json2 "encoding/json"
	"log"
	"bytes"
	"strings"
	"gopkg.in/kyokomi/emoji.v1"
)

/*
SendSkynetAlert will check if the message should be be, and if so build a telegram message
 */
func SendSkynetAlert(message string){
	if(!checkSend(message)){
		log.Println("Ignoreing message: "+message)
		return
	}

	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", message)
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprintf(":computer:"))
	buffer.WriteString(" ")
	buffer.WriteString("*Skynet Alert*\n")
	buffer.WriteString(dat["message"].(string))
	buffer.WriteString("\n")

	args := dat["args"].(map[string]interface{})
	config := args["config"].(map[string]interface{})
	chef_config := config["chef_config"].(map[string]interface{})
	buffer.WriteString("*Environment: *")
	buffer.WriteString(chef_config["environment"].(string))
	buffer.WriteString("\n")

	runlist := chef_config["run_list"].([]interface {})
	log.Println(runlist)
	buffer.WriteString("*Run list: *")

	for _, recipe := range runlist {
		recipeS := strings.Replace(recipe.(string),"recipe[","",-1)
		recipeS = strings.Replace(recipeS,"]","",-1)
		buffer.WriteString(recipeS)
		buffer.WriteString(" ")
	}

	buffer.WriteString("\n")

	cmdb_config := config["cmdb_config"].(map[string]interface{})

	buffer.WriteString("*Description: *")
	buffer.WriteString(cmdb_config["description"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Technical: *")
	buffer.WriteString(cmdb_config["technical"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Environment: *")
	buffer.WriteString(cmdb_config["environment"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Application: *")
	buffer.WriteString(cmdb_config["application"].(string))
	buffer.WriteString("\n")

	SendAlert(buffer.String())

}
