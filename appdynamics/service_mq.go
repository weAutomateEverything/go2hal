package appdynamics

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type MqService interface {
	addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string, ignorePrefix []string) error
}

type mqService struct {
	alert alert.Service
	store Store
}

func NewMqSercvice(alertService alert.Service, store Store) MqService {
	go func() {
		monitorAppdynamicsQueue(store, alertService)
	}()
	return &mqService{
		store: store,
		alert: alertService,
	}
}

func (s *mqService) addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string, ignorePrefix []string) error {
	endpointObject := MqEndpoint{MetricPath: metricPath, Application: application, Name: name}
	err := checkQueues(endpointObject, s.alert, s.store, chatId)
	if err != nil {
		return err
	}
	s.store.addMqEndpoint(name, application, metricPath, chatId, ignorePrefix)
	return nil
}

func monitorAppdynamicsQueue(s Store, a alert.Service) {
	for true {
		endpoints, err := s.getAllEndpoints()
		if err != nil {
			log.Printf("Error on DB MQ query: %s", err.Error())
		} else {
			for _, endpoint := range endpoints {
				for _, queue := range endpoint.MqEndpoints {
					if queue.Disabled == false {
						checkQueues(*queue, a, s, endpoint.ChatId)
					}
				}
			}
		}
		time.Sleep(time.Minute * 30)
	}
}

func checkQueues(endpoint MqEndpoint, a alert.Service, s Store, chat uint32) (err error) {
	ctx := context.Background()
	response, err := doGet(ctx, buildQueryString(endpoint), s, a, chat)
	if err != nil {
		log.Printf("%v, Unable to query appdynamics: %s", endpoint.Name, err.Error())
		a.SendError(ctx, fmt.Errorf("%v we are unable to query app dynamics", endpoint.Name))
		a.SendError(ctx, fmt.Errorf(" Queue depths Error: %s", err.Error()))
		return err
	}

	var dat []interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		a.SendError(ctx, fmt.Errorf("%v, error unmarshalling: %s", endpoint.Name, err))
		return err
	}

	success := false
	for _, queue := range dat {
		queue2 := queue.(map[string]interface{})
		err = checkQueue(ctx, endpoint, queue2["name"].(string), a, s, chat)
		if err == nil {
			success = true
		} else {
			log.Printf("MQ Error: %v", err.Error())
		}

	}
	if !success {
		return a.SendAlert(ctx, chat, fmt.Sprintf("%v, no queues had any data. please check the machine agents are sending the data", endpoint.Name))
	}
	return nil
}
func checkQueue(ctx context.Context, endpoint MqEndpoint, name string, a alert.Service, s Store, chat uint32) error {
	for _, x := range endpoint.IgnorePrefix {
		if strings.HasPrefix(name, x) {
			return nil
		}
	}
	currDepth, err := getAppdValue(ctx, buildQueryStringQueueDepth(endpoint, name), s, a, chat)
	if err != nil {
		return err
	}

	maxDepth, err := getAppdValue(ctx, buildQueryStringMaxQueueDepth(endpoint, name), s, a, chat)
	if err != nil {
		return err
	}

	if maxDepth == 0 {
		return fmt.Errorf("%v, max depth for queue %s is 0", endpoint.Name, name)
	}
	if currDepth == 0 {
		return nil
	}
	full := currDepth / maxDepth * 100
	if full > 90 {

		a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :interrobang: %s - Queue %s, is more than 90 percent full. "+
			"Current Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))

		return nil
	}

	if full > 75 {
		a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s, is more than 75 percent full. Current "+
			"Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))

		return nil
	}

	if strings.HasSuffix(strings.ToUpper(name), "BK") && currDepth > 0 {
		a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :warning: %s - Backout Queue %s contains data. The queue should be empty.\nPlease investigate why messages are being placed on this queue.\nCurrent "+
			"Depth %.0f, Max Depth %.0f", endpoint.Name, name, currDepth, maxDepth))
		return nil
	}

	messageAge, err := getAppdValue(ctx, buildQueryStringOldestMessageAge(endpoint, name), s, a, chat)

	if messageAge > endpoint.MaxMessageAge {
		t := time.Now().Add(time.Second * time.Duration(messageAge) * -1)
		d := time.Since(t).Truncate(time.Second)
		if currDepth == 0 {
			a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s contains messages that are %v old. Queue depth is 0, so it looks like there are uncommitted messages stuck in the application. Please investigate why messages are not being processed on the queue.", endpoint.Name, name, d.String()))
			return nil
		}
		a.SendAlert(ctx, chat, emoji.Sprintf(":baggage_claim: :warning: %s - Queue %s contains messages that are %v old. Please investigate why messages are not being processed on the queue.\nCurrent "+
			"Depth %.0f, Max Depth %.0f", endpoint.Name, name, d.String(), currDepth, maxDepth))
		return nil

	}
	return nil
}

func getAppdValue(ctx context.Context, path string, s Store, a alert.Service, chat uint32) (float64, error) {
	response, err := doGet(ctx, path, s, a, chat)
	if err != nil {
		log.Printf("Error retreiving queue %s", err)
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

func buildQueryStringOldestMessageAge(endpoint MqEndpoint, queue string) string {
	return fmt.Sprintf("/controller/rest/applications/%s/metric-data?metric-path=%s%%7C%s%%7COldestMsgAge&time-range-type=BEFORE_NOW&duration-in-mins=15&output=JSON", endpoint.Application, endpoint.MetricPath, queue)

}
func doGet(ctx context.Context, uri string, s Store, a alert.Service, chat uint32) (string, error) {

	appd, err := s.GetAppDynamics(chat)
	if err != nil {
		a.SendError(ctx, err)
		return "", err
	}

	url := appd.Endpoint + uri

	resp, err := http.Get(url)
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
