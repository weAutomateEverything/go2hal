package rest

import (
	"net/http"
	"github.com/zamedic/go2hal/database"
	"encoding/json"
	"log"
	"github.com/zamedic/go2hal/service"
)

func addSeleniumCheck(w http.ResponseWriter, r *http.Request) {

	var selenium database.Selenium
	err := json.NewDecoder(r.Body).Decode(&selenium)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = service.TestSelenium(selenium)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}


}
