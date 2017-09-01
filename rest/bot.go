package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
)

type bot struct {
	Token string
}

func addBot(w http.ResponseWriter, r *http.Request) {
	var botObject bot
	_ = json.NewDecoder(r.Body).Decode(&botObject)
	database.AddBot(botObject.Token)
	w.WriteHeader(http.StatusOK)
}

func botStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	heartbeat := database.GetBotHeartbeat()
	json.NewEncoder(w).Encode(&heartbeat)
}
