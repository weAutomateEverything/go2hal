package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
)

type jira struct {
	URL, Template, DefaultUser string
}

type callout struct {
	URL string
}

type seleniumTimout struct {
	Timeout int
}

func saveJira(w http.ResponseWriter, r *http.Request) {
	var j jira
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.SaveJiraDetails(j.URL, j.Template, j.DefaultUser)
	w.WriteHeader(http.StatusOK)
}

func saveCallout(w http.ResponseWriter, r *http.Request) {
	var c callout
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.SaveCalloutDetails(c.URL)
	w.WriteHeader(http.StatusOK)
}


func saveSeleniumTimout(w http.ResponseWriter, r *http.Request) {
	var s seleniumTimout
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.SaveSeleniumTimeDetails(s.Timeout)
	w.WriteHeader(http.StatusOK)
}