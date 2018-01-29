package sensu

import (
	"strings"
	"fmt"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/zamedic/go2hal/alert"
)

type Service interface {
	handleSensu(sensu sensuMessage)
}

type service struct {
	alert alert.Service
}

func NewService(service alert.Service) Service{
	return &service{alert:service}
}


func (s *service) handleSensu(sensu sensuMessage) {
	for _, msg := range sensu.Attachments {
		e := ""
		if strings.Index(msg.Title, "CRITICAL") > 0 {
			e = ":rotating_light:"

		} else if strings.Index(msg.Title, "WARNING") > 0 {
			e = ":warning:"
		} else {
			e = ":white_check_mark:"
		}
		msg := fmt.Sprintf("%v *%v*\n %v", e, msg.Title, msg.Text)
		s.alert.SendAlert(emoji.Sprint(msg))

	}

}
