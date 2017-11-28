package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"encoding/json"
)

type response struct {
	Bots   [] database.HeartBeat
	Values map[string]int64
}

func status(w http.ResponseWriter, r *http.Request) {

	res := response{}
	res.Bots = database.GetBotHeartbeat()
	stats := database.GetStats()
	res.Values = stats.Values

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)

}
