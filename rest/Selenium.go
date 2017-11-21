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
	log.Println(selenium)
	log.Println(selenium.Pages[0].Actions[0])
	log.Println(selenium.Pages[0].Actions[0].Selector)
	log.Println(selenium.Pages[0].Actions[0].InputData)
	log.Println(selenium.Pages[0].Actions[0].ClickButton)
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
