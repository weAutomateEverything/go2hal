package jira

import (
	"strings"
	"bytes"
	"io/ioutil"
	"fmt"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/zamedic/go2hal/config"
	"github.com/zamedic/go2hal/alert"
	"log"
	"text/template"
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/user"
)

type Service interface {
	CreateJira(title, description string,name string)
}

type service struct {
	alert alert.Service
	confStore config.Store
	userStore user.Store
}


func (s *service)CreateJira(title, description string, username string) {
	title = strings.Replace(title, "\n", "", -1)
	description = strings.Replace(description, "\n", "", -1)
	type q struct {
		User        string
		Title       string
		Description string
	}

	j, err := s.confStore.GetJiraDetails()
	if err != nil {
		s.alert.SendError(err)
		return
	}
	if j == nil || j.URL == "" {
		log.Println("No JIRA URL Set. Will not create a JIRA Item")
		return
	}

	qr := q{User: jiraUser(username), Description: description, Title: title}
	tmpl, err := template.New("jira").Parse(j.Template)
	if err != nil {
		s.alert.SendError(err)
		return
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, qr)
	if err != nil {
		s.alert.SendError(err)
		return
	}

	resp, err := http.Post(j.URL, "application/json", buf)
	if err != nil {
		s.alert.SendError(err)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.alert.SendError(err)
		return
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		s.alert.SendError(fmt.Errorf("JIRA Response failed to unmarshall: %s, %s", err, response))
		return
	}

	key := dat["key"].(string)

	s.alert.SendAlert(emoji.Sprintf(":ticket: JIRA Ticket Created. ID: %s assigned to %s. Description %s", key, qr.User, description))
}


func jiraUser(username string, s service) string {
	j, err := s.confStore.GetJiraDetails()
	if j == nil {
		return ""
	}
	if err != nil {
		s.alert.SendError(err)
		return ""
	}
	if username == "DEFAULT" || err != nil {
		return j.DefaultUser
	}
	u := s.userStore.FindUserByCalloutName(username).JIRAName
	if u == "" {
		return j.DefaultUser
	}
	return u
}