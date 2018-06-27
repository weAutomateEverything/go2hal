package prometheus

import (
	"encoding/json"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
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
		msg = msg + a["labels"].(string) + "\n"
		msg = msg + a["annotations"].(string)
		s.alertService.SendAlert(context.TODO(), chat, msg)
	}

	return nil
}
