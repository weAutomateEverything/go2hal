package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/user"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Service interface {
	CreateJira(ctx context.Context, title, description string, name string)
}

type service struct {
	alert     alert.Service
	userStore user.Store
}

func NewService(alert alert.Service, userStore user.Store) Service {
	return &service{alert, userStore}

}

func (s *service) CreateJira(ctx context.Context, title, description string, username string) {
	title = strings.Replace(title, "\n", "", -1)
	description = strings.Replace(description, "\n", "", -1)
	title = strings.Replace(title, "\t", "", -1)
	description = strings.Replace(description, "\t", "", -1)
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
	description = strings.Replace(description, "\n", "\\n", -1)

	qr := q{User: s.jiraUser(username), Description: description, Title: title}
	tmpl, err := template.New("jira").Parse(jiraTemplate)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, qr)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}

	resp, err := http.Post(jiraUrl, "application/json", buf)
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	log.Printf("JIRA Response: %v", string(response))
	if err != nil {
		s.alert.SendError(ctx, err)
		return
	}

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(response), &dat); err != nil {
		s.alert.SendError(ctx, fmt.Errorf("JIRA Response failed to unmarshall: %s, %s", err, response))
		return
	}

	key, ok := dat["key"].(string)

	if !ok {
		s.alert.SendError(ctx, fmt.Errorf("JIRA Response dat key not a string: %s", response))
	}

	s.alert.SendAlert(ctx, emoji.Sprintf(":ticket: JIRA Ticket Created. ID: %s assigned to %s. Description %s", key, qr.User, description))
}

func (s *service) jiraUser(username string) string {
	if username == "DEFAULT" {
		return os.Getenv("JIRA_DEFAULT_USER")
	}
	u := s.userStore.FindUserByCalloutName(username).JIRAName
	if u == "" {
		return os.Getenv("JIRA_DEFAULT_USER")
	}
	return u
}
