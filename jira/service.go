package jira

import (
	"strings"
	"bytes"
	"io/ioutil"
	"fmt"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/zamedic/go2hal/alert"
	"log"
	"text/template"
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/user"
	"os"
)

type Service interface {
	CreateJira(title, description string,name string)
}

type service struct {
	alert alert.Service
	userStore user.Store
}

func NewService(alert alert.Service,userStore user.Store) Service{
	return &service{alert,userStore}

}

func (s *service)CreateJira(title, description string, username string) {
	title = strings.Replace(title, "\n", "", -1)
	description = strings.Replace(description, "\n", "", -1)
	type q struct {
		User        string
		Title       string
		Description string
	}

	jiraUrl := os.Getenv("JIRA_URL")
	if jiraUrl == "" {
		log.Println("No JIRA URL Set. Will not create a JIRA Item")
		return
	}

	jiraTemplate := os.Getenv("JIRA_TEMPLATE")

	qr := q{User: s.jiraUser(username), Description: description, Title: title}
	tmpl, err := template.New("jira").Parse(jiraTemplate)
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

	resp, err := http.Post(jiraUrl, "application/json", buf)
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


func (s *service)jiraUser(username string) string {
	if username == "DEFAULT"  {
		return os.Getenv("JIRA_DEFAULT_USER")
	}
	u := s.userStore.FindUserByCalloutName(username).JIRAName
	if u == "" {
		return os.Getenv("JIRA_DEFAULT_USER")
	}
	return u
}