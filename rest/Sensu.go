package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/service"
	"encoding/json"
	"log"
	"gopkg.in/kyokomi/emoji.v1"
	"strings"
	"fmt"
)

type sensuMessage struct {
	Text        string            `json:"text"`
	IconURL     string            `json:"icon_url"`
	Attachments []sensuAttachment `json:"attachments"`
}

type sensuAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func sensuSlackAlert(w http.ResponseWriter, r *http.Request) {
	var sensu sensuMessage

	err := json.NewDecoder(r.Body).Decode(&sensu)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, msg := range sensu.Attachments {
		e := ""
		if strings.Index(msg.Title, "CRITICAL") > 0 {
			e = ":rotating_light:"

		} else if strings.Index(msg.Title, "WARNING") > 0 {
			e = ":warning:"
		} else {
			e = ":white_check_mark:"
		}
		s := fmt.Sprintf("%v *%v*\n %v", e,msg.Title, msg.Text)
		service.SendAlert(emoji.Sprint(s))

	}

}
