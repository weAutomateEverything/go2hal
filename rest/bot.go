package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/telegram"
)

type bot struct {
	Token string
}

func addBot(w http.ResponseWriter, r *http.Request) {
	var botObject bot
	_ = json.NewDecoder(r.Body).Decode(&botObject)
	if botObject.Token == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bot Token cannot be empty"))
		return
	}

	err := telegram.TestBot(botObject.Token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	database.AddBot(botObject.Token)
	w.WriteHeader(http.StatusOK)
}

func botStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	heartbeat := database.GetBotHeartbeat()
	json.NewEncoder(w).Encode(&heartbeat)
}
