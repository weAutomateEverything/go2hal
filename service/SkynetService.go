package service

import (
	"log"
	"fmt"
	"net/http"
	"github.com/zamedic/go2hal/database"
	"io/ioutil"
	"strings"
	json2 "encoding/json"
	"time"
	"gopkg.in/kyokomi/emoji.v1"
	"bytes"
	"gopkg.in/telegram-bot-api.v4"
)

func init(){
	log.Println("Initializing Skynet Rebuild Command")
	register(func() command {
		return &rebuildNode{}
	})
}

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
		SendError(fmt.Errorf("skynet alert - rrror unmarshalling: %s", message))
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
	chefConfig := config["chef_config"].(map[string]interface{})
	buffer.WriteString("*Environment: *")
	buffer.WriteString(chefConfig["environment"].(string))
	buffer.WriteString("\n")

	runlist := chefConfig["run_list"].([]interface {})
	log.Println(runlist)
	buffer.WriteString("*Run list: *")

	for _, recipe := range runlist {
		recipeS := strings.Replace(recipe.(string),"recipe[","",-1)
		recipeS = strings.Replace(recipeS,"]","",-1)
		buffer.WriteString(recipeS)
		buffer.WriteString(" ")
	}

	buffer.WriteString("\n")

	cmdbConfig := config["cmdb_config"].(map[string]interface{})

	buffer.WriteString("*Description: *")
	buffer.WriteString(cmdbConfig["description"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Technical: *")
	buffer.WriteString(cmdbConfig["technical"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Environment: *")
	buffer.WriteString(cmdbConfig["environment"].(string))
	buffer.WriteString("\n")

	buffer.WriteString("*Application: *")
	buffer.WriteString(cmdbConfig["application"].(string))
	buffer.WriteString("\n")

	SendAlert(buffer.String())

}

/*
RecreateNode will find a node, get the details, delete the node, then recreate it
 */
func RecreateNode(nodeName, callerName string) error {
	skynet, err := database.GetSkynetRecord()
	if err != nil {
		return err
	}
	json, err := findNode(nodeName, skynet)
	if err != nil {
		return err
	}

	err = deleteNode(nodeName, callerName, skynet)
	if err != nil {
		return err
	}

	err = waitForDelete(nodeName, skynet)
	if err != nil {
		return err
	}

	err = createNode(json, skynet)
	err = waitForBuild(nodeName,skynet)
	return nil

}
func findNode(nodeName string, skynet database.Skynet) (string, error) {
	body, err := doHTTP("GET", skynet.Address+"/virtual_machines/"+nodeName, "", skynet)
	if err != nil {
		return "", err
	}
	return body, nil
}

func deleteNode(nodeName, callerName string, skynet database.Skynet) error {
	SendAlert(fmt.Sprintf("Received a Delete Node request from %s for node %s. Proceeding with Delete", callerName, nodeName))
	_, err := doHTTP("DELETE", skynet.Address+"/virtual_machines/"+nodeName, "", skynet)
	if err != nil {
		return err
	}
	return nil
}

func waitForDelete(nodeName string, skynet database.Skynet) error {
	return poll("ARCHIVED", nodeName, skynet, true)
}

func createNode(json string, skynet database.Skynet) error {

	body, err := doHTTP("POST", skynet.Address+"/virtual_machines", json, skynet)
	if err != nil {
		SendAlert(fmt.Sprintf("Error creating node. %s", err.Error()))
		return err
	}
	log.Println(body)
	return nil
}

func waitForBuild(nodeName string, skynet database.Skynet) error {
	return poll("BOOTSTRAPPED", nodeName, skynet, false)
}

func doHTTP(method, url, body string, skynet database.Skynet) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		logError(fmt.Sprintf("Error creating URL to find node %s", err.Error()))
		return "", err
	}
	req.SetBasicAuth(skynet.Username, skynet.Password)
	req.Header.Set("Content-Type","application/json")
	resp, err := client.Do(req)
	if err != nil {
		logError(fmt.Sprintf("error in skynet call. %s", err.Error()))
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	return s, nil
}

func logError(error string) {
	log.Println(error)
	SendAlert(error)
}

func poll(expectedState, nodeName string, skynet database.Skynet, ignoreFailed bool) error {
	i := 0
	for i < 300 {
		body, err := doHTTP ("GET", skynet.Address+"/virtual_machines/"+nodeName+"/state", "", skynet)
		if err != nil {
			logError(fmt.Sprintf("Error retreiving node state: %s", err.Error()))
			return err
		}
		var dat map[string]interface{}
		err = json2.Unmarshal([]byte(body), &dat)
		if err != nil {
			logError(fmt.Sprintf("error unamrhsaling %s to json", body))
			return err
		}
		state := dat["state"].(map[string]interface{})["current"].(string)
		if strings.ToUpper(state) == expectedState {
			SendAlert(fmt.Sprintf("%s has been reached state %s.", nodeName, expectedState))
			return nil
		}
		if (!ignoreFailed && strings.ToUpper(state) == "FAILED") {
			SendAlert(fmt.Sprintf("%s has entered a Failed State.", nodeName))
			return fmt.Errorf("%s has entered a Failed State", nodeName)
		}
		i++
		if i%30 == 0 {
			SendAlert(fmt.Sprintf("still waiting for node %s to reach state %s. Curent state is %s", nodeName,expectedState, state))
		}
		time.Sleep(time.Second)
	}
	err := fmt.Errorf("Timed out waiting for node %s to delete", nodeName)
	logError(err.Error())
	return err
}



type rebuildNode struct {

}

func (s *rebuildNode) commandIdentifier() string {
	return "RebuildNode"
}

func (s *rebuildNode) commandDescription() string {
	return "Rebuilds a node"
}

func (s *rebuildNode) execute(update tgbotapi.Update){
	RecreateNode(update.Message.CommandArguments(),update.Message.From.UserName)
}

