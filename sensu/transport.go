package sensu

import (
	"context"
	"net/http"
	"encoding/json"
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



func decodeSensu(_ context.Context, r *http.Request) (interface{}, error) {
	var sensu sensuMessage
	err := json.NewDecoder(r.Body).Decode(&sensu)
	return sensu, err
}
