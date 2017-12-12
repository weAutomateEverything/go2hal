package service

import (
	"log"
	"time"
	"github.com/zamedic/go2hal/database"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	json2 "encoding/json"
	"gopkg.in/kyokomi/emoji.v1"
	"strings"
	"errors"
)

func init() {
	go func() {
		monitorAppdynamicsQueue()
	}()
}

/*
SendAppdynamicsAlert will send an alert message from the data contained in the even message field
 */
func SendAppdynamicsAlert(message string) {

	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		SendError(fmt.Errorf("error unmarshalling App Dynamics Message: %s", message))
		return
	}

	events := dat["events"].([]interface{})
	for _, event := range events {
		event := event.(map[string]interface{})

		message := event["eventMessage"].(string)

		application := event["application"].(map[string]interface{})
		tier := event["tier"].(map[string]interface{})
		node := event["node"].(map[string]interface{})

		message = strings.Replace(message, "<b>", "*", -1)
		message = strings.Replace(message, "</b>", "*", -1)
		message = strings.Replace(message, "<br>", "\n", -1)

		var buffer bytes.Buffer
		buffer.WriteString(emoji.Sprintf(":red_circle:"))
		buffer.WriteString(" ")
		buffer.WriteString(message)
		buffer.WriteString("\n")

		if application != nil {
			app := application["name"].(string)
			app = strings.Replace(app, "_", "\\_", -1)
			buffer.WriteString("*Application:* ")
			buffer.WriteString(app)
			buffer.WriteString("\n")
		}

		if tier != nil {
			ti := tier["name"].(string)
			ti = strings.Replace(ti, "_", "\\_", -1)
			buffer.WriteString("*Tier:* ")
			buffer.WriteString(ti)
			buffer.WriteString("\n")
		}

		if node != nil {
			no := node["name"].(string)
			no = strings.Replace(no, "_", "\\_", -1)
			buffer.WriteString("*Node:* ")
			buffer.WriteString(no)
			buffer.WriteString("\n")
		}

		log.Printf("Sending Alert %s", buffer.String())
		SendAlert(buffer.String())
	}
}

func SendNonTechnicalAppDynamicsAlert(message string){

	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprintf(":red_circle:"))
	buffer.WriteString(" ")
	buffer.WriteString(message)

	log.Printf("Sending Alert %s", message)
	SendAlert(buffer.String())
}

/*
AddAppdynamicsEndpoint calls the mongo service
 */
func AddAppdynamicsEndpoint(endpoint string) error {
	return database.AddAppDynamicsEndpoint(endpoint)
}

/*
AddAppDynamicsQueue validated that the input data is valid, then persists the data to the mongo database.
 */
func AddAppDynamicsQueue(name, application, metricPath string) error {
	endpointObject := database.MqEndpoint{MetricPath: metricPath, Application: application, Name: name};
	err := checkQueues(endpointObject);
	if err != nil {
		return err
	}
	database.AddMqEndpoint(name, application, metricPath)
	return nil;
}

/*
ExecuteCommandFromAppd will search for the nodes IP Address, then execute the command on the node
*/
func ExecuteCommandFromAppd(commandName, applicationID, nodeId string) error {
	ipaddress, err := getIPAddressForNode(applicationID, nodeId)
	if err != nil {
		SendError(err)
		return err
	}
	return ExecuteRemoteCommand(commandName,ipaddress)
}

func monitorAppdynamicsQueue() {
	log.Println("Starting App Dynamics Queue Service")
	for true {
		endpoints, err := database.GetAppDynamics()
		if err != nil {
			log.Printf("Error on MQ query: %s", err.Error())
		} else {
			for _, endpoint := range endpoints.MqEndpoints {
				checkQueues(endpoint)
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func checkQueues(endpoint database.MqEndpoint) error {

	response, err := doGet(buildQueryString(endpoint))
	if err != nil {
		log.Printf("Unable to query appdynamics: %s", err.Error())
		SendError(fmt.Errorf("we are unable to query app dynamics"))
		SendError(fmt.Errorf(" Queue depths Error: %s", err.Error()))
		return err;
	}

	if err != nil {
		SendError(fmt.Errorf("queue - error parsing body %s", err))
		return err
	}

	var dat []interface{}
	if err := json2.Unmarshal([]byte(response), &dat); err != nil {
		SendError(fmt.Errorf("error unmarshalling: %s", err))
		return err
	}

	success := false
	for _, queue := range dat {
		queue2 := queue.(map[string]interface{})
		err = checkQueue(endpoint, queue2["name"].(string))
		if err == nil {
			success = true
		}

	}
	if !success {
		SendError(errors.New("no queues had any data. please check the machine agents are sending the data"))
	}
	return nil
}

func checkQueue(endpoint database.MqEndpoint, name string) error {
	currDepth, err := getCurrentQueueDepthValue(buildQueryStringQueueDepth(endpoint, name))
	if err != nil {
		return err
	}

	maxDepth, err := getCurrentQueueDepthValue(buildQueryStringMaxQueueDepth(endpoint, name))
	if err != nil {
		return err
	}

	if maxDepth == 0 {
		return fmt.Errorf("max depth for queue %s is 0", name)
	}
	full := currDepth / maxDepth * 100;
	if full > 90 {
		SendAlert(emoji.Sprintf(":baggage_claim: :interrobang: %s - Queue %s, is more than 90 percent full. "+
			"Current Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		return nil
	}

	if full > 75 {
		SendAlert(emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s, is more than 75 percent full. Current "+
			"Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		return nil
	}
	return nil
}

func getCurrentQueueDepthValue(path string) (float64, error) {
	response, err := doGet(path)
	if err != nil {
		log.Printf("Error retreiving queue %s", err)
		return 0, err
	}
	if err != nil {
		log.Printf("Error reading body %s", err)
		return 0, err
	}

	var dat []interface{}
	if err := json2.Unmarshal([]byte(response), &dat); err != nil {
		log.Printf("error unmarshalling body %s", err)
		return 0, err
	}

	if len(dat) == 0 {
		return 0, fmt.Errorf("no data found for %s", path)
	}

	record := dat[0].(map[string]interface{})
	values := record["metricValues"].([]interface{})
	if len(values) == 0 {
		return 0, nil;
	}
	value := values[0].(map[string]interface{})
	return value["current"].(float64), nil

}

func buildQueryString(endpoint database.MqEndpoint) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metrics?metric-path=%s&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath)
}

func buildQueryStringQueueDepth(endpoint database.MqEndpoint, queue string) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metric-data?metric-path=%s%%7C%s%%7CCurrent%%20Queue%%20Depth&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath, queue)
}

func buildQueryStringMaxQueueDepth(endpoint database.MqEndpoint, queue string) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metric-data?metric-path=%s%%7C%s%%7CMax%%20Queue%%20Depth&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath, queue)

}

func doGet(uri string) (string, error) {

	a, err := database.GetAppDynamics()
	if err != nil {
		SendError(err)
		return "", err
	}

	client := &http.Client{}
	url := a.Endpoint + uri
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		SendError(err)
		return "", err
	}

	//req.SetBasicAuth(a.User, a.Password) //<username@group>
	resp, err := client.Do(req)
	if err != nil {
		SendError(err)
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		SendError(err)
		return "", err
	}
	s := string(bodyText)
	return s, nil
}

func getIPAddressForNode(application, node string) (string, error) {
	uri := fmt.Sprintf("/controller/rest/applications/%s/nodes/%s?output=json", application, node)
	response, err := doGet(uri)
	log.Println(response)
	if err != nil {
		SendError(err)
		return "", err
	}

	var dat []interface{}
	err = json2.Unmarshal([]byte(response), &dat)
	if err != nil {
		SendError(err)
		return "",err
	}
	v := dat[0].(map[string]interface{})
	ipaddresses := v["ipAddresses"].(map[string]interface{})
	arrayIP := ipaddresses["ipAddresses"].([]interface{})
	for _,ipo := range arrayIP {
		ip := ipo.(string)
		if (strings.Index(ip,".") > 0){
			return ip,nil
		}
	}
	return "", errors.New("no up address found")
}
