package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"encoding/json"
)

type response struct {
	Bots                 [] database.HeartBeat
	MessagesSent         int64
	MessagesReceived     int64
	AppDynamicsMessages  int64
	ChefDeliveryMessages int64
	MessageQueueDepth    int
}

func status(w http.ResponseWriter, r *http.Request) {

	res := response{}
	res.Bots = database.GetBotHeartbeat()
	res.MessagesSent, res.MessagesReceived, res.AppDynamicsMessages, res.ChefDeliveryMessages = database.GetStats()
	res.MessageQueueDepth = database.QueueDepth()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)

}
