package prometheus

import (
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"sort"
)

type Service interface {
	sendPrometheusAlert(chat uint32, message string) error
}

func NewService(alertService alert.Service) Service {
	return &service{
		alertService: alertService,
	}
}

type service struct {
	alertService alert.Service
}

func (s *service) sendPrometheusAlert(chat uint32, message string) (err error) {
	var r map[string]interface{}

	err = json.Unmarshal([]byte(message), &r)
	if err != nil {
		return
	}

	alerts := r["alerts"].([]interface{})
	for _, alert := range alerts {
		a := alert.(map[string]interface{})
		msg := a["status"].(string) + "\n"
		labels := a["labels"].(map[string]interface{})
		keys := make([]string, 1)
		for key, _ := range labels {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			if key != "" {
				msg = msg + fmt.Sprintf("*%v*: %v\n", key, labels[key])
			}
		}
		annotation := a["annotations"].(map[string]interface{})
		for key, value := range annotation {
			msg = msg + fmt.Sprintf("*%v*: %v\n", key, value)
		}
		s.alertService.SendAlert(context.TODO(), chat, msg)
	}

	return nil
}
