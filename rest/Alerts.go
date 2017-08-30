package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"fmt"
	"github.com/zamedic/go2hal/telegram"
	"io/ioutil"
)

func AlertHandler(w http.ResponseWriter, r *http.Request){
	alertGroup,err := database.AlertGroup()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w,err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	telegram.SendMessage(alertGroup,string(body),0)
	w.WriteHeader(http.StatusOK)
}
