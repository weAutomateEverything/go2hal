package appdynamics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/ssh"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	sendAppdynamicsAlert(ctx context.Context, message string)
	addAppdynamicsEndpoint(endpoint string) error
	addAppDynamicsQueue(name, application, metricPath string) error
	executeCommandFromAppd(ctx context.Context, commandName, applicationID, nodeID string) error
}

type service struct {
	alert alert.Service
	store Store
	ssh   ssh.Service
}

func NewService(alertService alert.Service, sshservice ssh.Service, store Store) Service {
	go func() {
		monitorAppdynamicsQueue(store, alertService)
	}()
	return &service{alert: alertService, ssh: sshservice, store: store}
}

func (s *service) sendAppdynamicsAlert(ctx context.Context, message string) {

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(message), &dat); err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error unmarshalling App Dynamics Message: %s", message))
		return
	}

	if val, ok := dat["business"]; ok {
		business := val.(map[string]interface{})

		var NonTechBuffer bytes.Buffer

		businessMessage := business["businessEvent"].(string)
		NonTechBuffer.WriteString(emoji.Sprintf(":red_circle:"))
		NonTechBuffer.WriteString(" ")
		NonTechBuffer.WriteString(businessMessage)

		log.Printf("Sending Non-Technical Alert %s", NonTechBuffer.String())
		s.alert.SendNonTechnicalAlert(ctx, NonTechBuffer.String())
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
		s.alert.SendAlert(ctx, buffer.String())
	}
}

func (s *service) addAppdynamicsEndpoint(endpoint string) error {
	return s.store.addAppDynamicsEndpoint(endpoint)
}

func (s *service) addAppDynamicsQueue(name, application, metricPath string) error {
	endpointObject := MqEndpoint{MetricPath: metricPath, Application: application, Name: name}
	err := checkQueues(endpointObject, s.alert, s.store)
	if err != nil {
		return err
	}
	s.store.addMqEndpoint(name, application, metricPath)
	return nil
}

func (s *service) executeCommandFromAppd(ctx context.Context, commandName, applicationID, nodeID string) error {
	ipaddress, err := s.getIPAddressForNode(applicationID, nodeID)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}
	return s.ssh.ExecuteRemoteCommand(ctx, commandName, ipaddress)
}

func monitorAppdynamicsQueue(s Store, a alert.Service) {
	for true {
		endpoints, err := s.GetAppDynamics()
		if err != nil {
			log.Printf("Error on MQ query: %s", err.Error())
		} else {
			for _, endpoint := range endpoints.MqEndpoints {
				checkQueues(endpoint, a, s)
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func checkQueues(endpoint MqEndpoint, a alert.Service, s Store) error {

	response, err := doGet(buildQueryString(endpoint), s, a)
	if err != nil {
		log.Printf("Unable to query appdynamics: %s", err.Error())
		a.SendError(context.TODO(), fmt.Errorf("we are unable to query app dynamics"))
		a.SendError(context.TODO(), fmt.Errorf(" Queue depths Error: %s", err.Error()))
		return err
	}

	if err != nil {
		a.SendError(context.TODO(), fmt.Errorf("queue - error parsing body %s", err))
		return err
	}

	var dat []interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		a.SendError(context.TODO(), fmt.Errorf("error unmarshalling: %s", err))
		return err
	}

	success := false
	for _, queue := range dat {
		queue2 := queue.(map[string]interface{})
		err = checkQueue(endpoint, queue2["name"].(string), a, s)
		if err == nil {
			success = true
		}

	}
	if !success {
		a.SendError(context.TODO(), errors.New("no queues had any data. please check the machine agents are sending the data"))
	}
	return nil
}
func checkQueue(endpoint MqEndpoint, name string, a alert.Service, s Store) error {
	currDepth, err := getCurrentQueueDepthValue(buildQueryStringQueueDepth(endpoint, name), s, a)
	if err != nil {
		return err
	}

	maxDepth, err := getCurrentQueueDepthValue(buildQueryStringMaxQueueDepth(endpoint, name), s, a)
	if err != nil {
		return err
	}

	if maxDepth == 0 {
		return fmt.Errorf("max depth for queue %s is 0", name)
	}
	full := currDepth / maxDepth * 100
	if full > 90 {
		a.SendAlert(context.TODO(), emoji.Sprintf(":baggage_claim: :interrobang: %s - Queue %s, is more than 90 percent full. "+
			"Current Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		return nil
	}

	if full > 75 {
		a.SendAlert(context.TODO(), emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s, is more than 75 percent full. Current "+
			"Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		return nil
	}
	return nil
}

func getCurrentQueueDepthValue(path string, s Store, a alert.Service) (float64, error) {
	response, err := doGet(path, s, a)
	if err != nil {
		log.Printf("Error retreiving queue %s", err)
		return 0, err
	}
	if err != nil {
		log.Printf("Error reading body %s", err)
		return 0, err
	}

	var dat []interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		log.Printf("error unmarshalling body %s", err)
		return 0, err
	}

	if len(dat) == 0 {
		return 0, fmt.Errorf("no data found for %s", path)
	}

	record := dat[0].(map[string]interface{})
	values := record["metricValues"].([]interface{})
	if len(values) == 0 {
		return 0, nil
	}
	value := values[0].(map[string]interface{})
	return value["current"].(float64), nil

}

func buildQueryString(endpoint MqEndpoint) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metrics?metric-path=%s&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath)
}

func buildQueryStringQueueDepth(endpoint MqEndpoint, queue string) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metric-data?metric-path=%s%%7C%s%%7CCurrent%%20Queue%%20Depth&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath, queue)
}

func buildQueryStringMaxQueueDepth(endpoint MqEndpoint, queue string) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metric-data?metric-path=%s%%7C%s%%7CMax%%20Queue%%20Depth&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath, queue)

}

func doGet(uri string, s Store, a alert.Service) (string, error) {

	appd, err := s.GetAppDynamics()
	if err != nil {
		a.SendError(context.TODO(), err)
		return "", err
	}

	client := &http.Client{}
	url := appd.Endpoint + uri
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		a.SendError(context.TODO(), err)
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		a.SendError(context.TODO(), err)
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		a.SendError(context.TODO(), err)
		return "", err
	}
	res := string(bodyText)
	return res, nil
}

func (s *service) getIPAddressForNode(application, node string) (string, error) {
	uri := fmt.Sprintf("/controller/rest/applications/%s/nodes/%s?output=json", application, node)
	response, err := doGet(uri, s.store, s.alert)
	log.Println(response)
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return "", err
	}

	var dat []interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return "", err
	}
	v := dat[0].(map[string]interface{})
	ipaddresses := v["ipAddresses"].(map[string]interface{})
	arrayIP := ipaddresses["ipAddresses"].([]interface{})
	for _, ipo := range arrayIP {
		ip := ipo.(string)
		if strings.Index(ip, ".") > 0 {
			return ip, nil
		}
	}
	return "", errors.New("no up address found")
}
