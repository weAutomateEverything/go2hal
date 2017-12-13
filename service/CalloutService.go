package service

import (
	"log"
	"os"
	"strconv"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/zamedic/go2hal/database"
	"html/template"
	"bytes"
	"fmt"
	json2 "encoding/json"
	"gopkg.in/kyokomi/emoji.v1"
	"errors"
	"runtime/debug"
)

/*
InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
 */
func InvokeCallout(title, message string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			SendError(errors.New(fmt.Sprint(err)))
			SendError(errors.New(string(debug.Stack())))

		}
	}()
	SendError(fmt.Errorf("invoking callout for: %s, %s",title,message))
	sendSNMPMessage()
	SendError(errors.New("Checking JIRA"))
	createJira(title, message)
}



func createJira(title, description string) {
	title = strings.Replace(title,"\n","",-1)
	description = strings.Replace(description,"\n","",-1)
	type q struct {
		User        string
		Title       string
		Description string
	}

	j, err := database.GetJiraDetails()
	if err != nil {
		SendError(err)
		return
	}
	if j == nil || j.URL == "" {
		log.Println("No JIRA URL Set. Will not create a JIRA Item")
		return
	}

	qr := q{User: jiraUser(), Description: description, Title:title}
	tmpl, err := template.New("jira").Parse(j.Template)
	if err != nil {
		SendError(err)
		return
	}
	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, qr)
	if err != nil {
		SendError(err)
		return
	}

	SendError(errors.New(buf.String()))

	resp, err := http.Post(j.URL, "application/json", buf)
	if err != nil {
		SendError(err)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendError(err)
		return
	}

	var dat map[string]interface{}
	if err := json2.Unmarshal([]byte(response), &dat); err != nil {
		SendError(fmt.Errorf("JIRA Response failed to unmarshall: %s, %s", err, response))
		return
	}

	key := dat["key"].(string)

	SendAlert(emoji.Sprintf(":ticket: JIRA Ticket Created. ID: %s assigned to %s. Description %s", key, qr.User, description))
}

func jiraUser() string {
	user, err := getFirstCallName()
	if err != nil {
		SendError(err)
		return ""

	}
	j, err := database.GetJiraDetails()
	if j == nil {
		return ""
	}
	if err != nil {
		SendError(err)
		return ""
	}
	if user == "DEFAULT" || err != nil {
		return j.DefaultUser
	}
	u := database.FindUserByCalloutName(user).JIRAName
	if u == "" {
		return j.DefaultUser
	}
	return u
}

func getFirstCallName() (string, error) {
	c, err := database.GetCalloutDetails()
	if err != nil {
		return "", err
	}
	if c == nil {
		return "", errors.New("no callout set")
	}
	endpoint := c.URL
	if endpoint == "" {
		return "DEFAULT", nil
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		SendError(err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		SendError(err)
		return "", err
	}
	bodyString := string(body)
	split := strings.SplitAfter(bodyString, "<font color='red' size=2>")
	names := strings.Split(split[1], "</font>")
	return names[0], nil

}

func snmpServier() string {
	return os.Getenv("SNMP_SERVER")
}

func snmpPort() uint16 {
	i, _ := strconv.ParseInt(os.Getenv("SNMP_PORT"), 10, 16)
	return uint16(i)
}
