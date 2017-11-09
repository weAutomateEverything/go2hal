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
)

func init() {
	go func() {
		monitorAppdynamicsQueue()
	}()
}

/*
AddAppDynamicsQueue validated that the input data is valid, then persists the data to the mongo database.
 */
func AddAppDynamicsQueue(name, endpoint, metricPath string) error {
	endpointObject := database.MqEndpoint{MetricPath: metricPath, Endpoint: endpoint, Name: name};
	err := checkQueues(endpointObject);
	if err != nil {
		return err
	}
	database.AddMqEndpoint(name, endpoint, metricPath)
	return nil;
}

func monitorAppdynamicsQueue() {
	log.Println("Starting App Dynamics Queue Service")
	for true {
		endpoints := database.GetMQEndponts()
		if endpoints != nil {
			for _, endpoint := range endpoints {
				checkQueues(endpoint)
			}

		}
		time.Sleep(time.Minute * 10)
	}
}
func checkQueues(endpoint database.MqEndpoint) error {
	log.Printf("Checking App Dynamics")

	response, err := http.Get(buildQueryString(endpoint))
	if err != nil {
		log.Printf("Unable to query appdynamics: %s",err.Error())
		SendAlert("We are unable to query App Dynamics to monitor Queue depths")
		return err;
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	log.Printf("Received: %s", body)

	if err != nil {
		log.Printf("Error parsing body %s", err)
		return err
	}

	var dat []interface{}
	if err := json2.Unmarshal([]byte(body), &dat); err != nil {
		log.Printf("Error unmarshalling: %s", err)
		return err
	}

	for _, queue := range dat {
		queue2 := queue.(map[string]interface{})
		checkQueue(endpoint, queue2["name"].(string))
	}
	return nil
}

func checkQueue(endpoint database.MqEndpoint, name string) error {
	log.Printf("Checking Queue %s", name)
	currDepth, err := getCurrentQueueDepthValue(buildQueryStringQueueDepth(endpoint, name))
	if err != nil {
		return err;
	}

	maxDepth, err := getCurrentQueueDepthValue(buildQueryStringMaxQueueDepth(endpoint, name))
	if err != nil {
		return err;
	}

	log.Printf("Queue: %s, Current Depth: %.0f, Max Depth: %.0f", name, currDepth, maxDepth)
	if maxDepth == 0 {
		return nil
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
	response, err := http.Get(path)
	if err != nil {
		log.Printf("Error retreiving queue %s", err)
		return 0, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	log.Printf("Received: %s", body)

	if err != nil {
		log.Printf("Error reading body %s", err)
		return 0, err
	}

	var dat []interface{}
	if err := json2.Unmarshal([]byte(body), &dat); err != nil {
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
	var buffer bytes.Buffer
	buffer.WriteString(endpoint.Endpoint)

	buffer.WriteString("/metrics?metric-path=")
	buffer.WriteString(endpoint.MetricPath)
	buffer.WriteString("&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON")
	return buffer.String()
}

func buildQueryStringQueueDepth(endpoint database.MqEndpoint, queue string) string {
	var buffer bytes.Buffer
	buffer.WriteString(endpoint.Endpoint)
	buffer.WriteString("/metric-data?metric-path=")
	buffer.WriteString(endpoint.MetricPath)
	buffer.WriteString("%7C")
	buffer.WriteString(queue)
	buffer.WriteString("%7CCurrent%20Queue%20Depth")

	buffer.WriteString("&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON")
	return buffer.String()
}

func buildQueryStringMaxQueueDepth(endpoint database.MqEndpoint, queue string) string {
	var buffer bytes.Buffer
	buffer.WriteString(endpoint.Endpoint)
	buffer.WriteString("/metric-data?metric-path=")
	buffer.WriteString(endpoint.MetricPath)
	buffer.WriteString("%7C")
	buffer.WriteString(queue)
	buffer.WriteString("%7CMax%20Queue%20Depth")
	buffer.WriteString("&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON")
	return buffer.String()
}
