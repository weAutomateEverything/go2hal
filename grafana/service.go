package grafana

import (
	"encoding/json"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"strings"
)

type Service interface {
	sendGrafanaAlert(ctx context.Context, chat uint32, body string) error
}

func NewService(alertService alert.Service) Service {
	return &service{
		alertService: alertService,
	}
}

type service struct {
	alertService alert.Service
}

func (s *service) sendGrafanaAlert(ctx context.Context, chat uint32, body string) (err error) {
	var r map[string]interface{}
	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		return
	}

	msg := r["title"].(string)
	msg = strings.Replace(msg, "[", "*", -1)
	msg = strings.Replace(msg, "]", "*", -1)
	msg = msg + "\n" + r["message"].(string)

	return s.alertService.SendAlert(ctx, chat, msg)

}
