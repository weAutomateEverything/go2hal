package appdynamics

import (
	appd "appdynamics"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/appdynamics/util"
	"github.com/weAutomateEverything/go2hal/ssh"
	"golang.org/x/net/context/ctxhttp"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Service interface {
	sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) error
	addAppdynamicsEndpoint(chat uint32, endpoint string) error
	addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string) error
	executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error
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

func (s *service) sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) error {
	var m appdynamicsMessage

	err := json.Unmarshal([]byte(message), &m)

	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
		s.alert.SendError(ctx, fmt.Errorf("error unmarshalling App Dynamics Message: %s. error: %v", message, err))
		return err
	}

	for _, event := range m.Events {

		message := event.EventMessage
		message = strings.Replace(message, "<b>", "*", -1)
		message = strings.Replace(message, "</b>", "*", -1)
		message = strings.Replace(message, "<br>", "\n", -1)

		var buffer bytes.Buffer
		buffer.WriteString(emoji.Sprintf(":red_circle:"))
		buffer.WriteString(" ")
		buffer.WriteString(message)
		buffer.WriteString("\n")

		if event.Application.Name != "" {
			app := strings.Replace(event.Application.Name, "_", "\\_", -1)
			buffer.WriteString("*Application:* ")
			buffer.WriteString(app)
			buffer.WriteString("\n")
		}

		if event.Tier.Name != "" {
			ti := strings.Replace(event.Tier.Name, "_", "\\_", -1)
			buffer.WriteString("*Tier:* ")
			buffer.WriteString(ti)
			buffer.WriteString("\n")
		}

		if event.Node.Name != "" {
			no := strings.Replace(event.Node.Name, "_", "\\_", -1)
			buffer.WriteString("*Node:* ")
			buffer.WriteString(no)
			buffer.WriteString("\n")
		}

		log.Printf("Sending Alert %s", buffer.String())
		tmperr := s.alert.SendAlert(ctx, chatId, buffer.String())
		if tmperr != nil {
			err = tmperr
		}
	}
	return err
}

func (s *service) addAppdynamicsEndpoint(chat uint32, endpoint string) error {
	return s.store.addAppDynamicsEndpoint(chat, endpoint)
}

func (s *service) addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string) error {
	endpointObject := MqEndpoint{MetricPath: metricPath, Application: application, Name: name}
	err := checkQueues(endpointObject, s.alert, s.store, chatId)
	if err != nil {
		return err
	}
	s.store.addMqEndpoint(name, application, metricPath, chatId)
	return nil
}

func (s *service) executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error {
	ipaddress, err := s.getIPAddressForNode(ctx, applicationID, nodeID, chatId)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}
	return s.ssh.ExecuteRemoteCommand(ctx, chatId, commandName, ipaddress)
}

func monitorAppdynamicsQueue(s Store, a alert.Service) {
	for true {
		endpoints, err := s.getAllEndpoints()
		if err != nil {
			log.Printf("Error on MQ query: %s", err.Error())
		} else {
			for _, endpoint := range endpoints {
				for _, queue := range endpoint.MqEndpoints {
					checkQueues(queue, a, s, endpoint.ChatId)
				}
			}
		}
		time.Sleep(time.Minute * 10)
	}
}

func checkQueues(endpoint MqEndpoint, a alert.Service, s Store, chat uint32) (err error) {
	handle, ctx := util.Start("appdynamics.checkQueues", "")
	defer appd.EndBT(handle)

	response, err := doGet(ctx, buildQueryString(endpoint), s, a, chat)
	if err != nil {
		log.Printf("Unable to query appdynamics: %s", err.Error())
		util.AddErrorToAppDynamics(ctx, err)
		a.SendError(ctx, fmt.Errorf("we are unable to query app dynamics"))
		a.SendError(ctx, fmt.Errorf(" Queue depths Error: %s", err.Error()))
		return err
	}

	var dat []interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		a.SendError(ctx, fmt.Errorf("error unmarshalling: %s", err))
		util.AddErrorToAppDynamics(ctx, err)
		return err
	}

	success := false
	for _, queue := range dat {
		queue2 := queue.(map[string]interface{})
		err = checkQueue(ctx, endpoint, queue2["name"].(string), a, s, chat)
		if err == nil {
			success = true
		}

	}
	if !success {
		a.SendError(ctx, errors.New("no queues had any data. please check the machine agents are sending the data"))
	}
	return nil
}
func checkQueue(ctx context.Context, endpoint MqEndpoint, name string, a alert.Service, s Store, chat uint32) error {
	currDepth, err := getCurrentQueueDepthValue(ctx, buildQueryStringQueueDepth(endpoint, name), s, a, chat)
	if err != nil {
		return err
	}

	maxDepth, err := getCurrentQueueDepthValue(ctx, buildQueryStringMaxQueueDepth(endpoint, name), s, a, chat)
	if err != nil {
		return err
	}

	if maxDepth == 0 {
		return fmt.Errorf("max depth for queue %s is 0", name)
	}
	full := currDepth / maxDepth * 100
	if full > 90 {
		for _, chat := range endpoint.Chat {
			a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :interrobang: %s - Queue %s, is more than 90 percent full. "+
				"Current Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		}

		return nil
	}

	if full > 75 {
		for _, chat := range endpoint.Chat {
			a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s, is more than 75 percent full. Current "+
				"Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		}
		return nil
	}
	return nil
}

func getCurrentQueueDepthValue(ctx context.Context, path string, s Store, a alert.Service, chat uint32) (float64, error) {
	response, err := doGet(ctx, path, s, a, chat)
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

func doGet(ctx context.Context, uri string, s Store, a alert.Service, chat uint32) (string, error) {

	appd, err := s.GetAppDynamics(chat)
	if err != nil {
		a.SendError(ctx, err)
		return "", err
	}

	client := &http.Client{}
	url := appd.Endpoint + uri
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		a.SendError(ctx, err)
		return "", err
	}

	resp, err := ctxhttp.Do(ctx, client, req)
	if err != nil {
		a.SendError(ctx, err)
		return "", err
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		a.SendError(ctx, err)
		return "", err
	}
	res := string(bodyText)
	return res, nil
}

func (s *service) getIPAddressForNode(ctx context.Context, application, node string, chat uint32) (string, error) {
	uri := fmt.Sprintf("/controller/rest/applications/%s/nodes/%s?output=json", application, node)
	response, err := doGet(ctx, uri, s.store, s.alert, chat)
	log.Println(response)
	if err != nil {
		s.alert.SendError(ctx, err)
		return "", err
	}

	var dat []interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		s.alert.SendError(ctx, err)
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
