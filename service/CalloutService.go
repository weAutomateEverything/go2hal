package service

import (
	g "github.com/soniah/gosnmp"
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
)

/*
InvokeCallout will invoke snmp if configured, then create a jira ticket if configured.
 */
func InvokeCallout(message string){
	sendSNMPMessage()
	createJira(message)
}

func sendSNMPMessage() {
	if snmpServier() == "" {
		return
	}
	g.Default.Port = snmpPort()
	g.Default.Target = snmpServier()
	g.Default.Version = g.Version2c
	g.Default.Logger = log.New(os.Stdout, "", 0)

	log.Printf("SNMP Server: %s Port: %d", g.Default.Target, g.Default.Port)

	err := g.Default.Connect()
	if err != nil {
		log.Printf("Connect() err: %v", err)
		return
	}
	defer g.Default.Conn.Close()

	p := g.SnmpPDU{
		Name:  ".1.3.6.1.4.1.789.1.2.2.4.0",
		Value: "Test Alert Message from HAL BOT. Please invoke Callout Group XXXXXXXXX",
		Type:  g.OctetString,
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{p},
	}

	result, err := g.Default.SendTrap(trap)
	if err != nil {
		log.Printf("Connect() err: %v", err)
		return
	}

	log.Printf("Error: %d", result.Error)
	log.Printf("Request ID %d", result.RequestID)
	SendAlert(emoji.Sprint(":telephone_receiver: Invoked callout"))

}

func createJira(description string) {
	type q struct {
		User        string
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

	qr := q{User: jiraUser(), Description: description}
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
	if err != nil{
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
