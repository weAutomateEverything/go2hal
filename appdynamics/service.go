package appdynamics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/ssh"
	"log"
	"strings"
)

type Service interface {
	sendAppdynamicsAlert(ctx context.Context, chatId uint32, message AppdynamicsMessage) error
	addAppdynamicsEndpoint(chat uint32, endpoint string) error
	executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error
}

type service struct {
	alert          alert.Service
	store          Store
	calloutService callout.Service
	ssh            ssh.Service
}

func NewService(alertService alert.Service, sshservice ssh.Service, store Store, callout callout.Service) Service {
	return &service{alert: alertService, ssh: sshservice, store: store, calloutService: callout}
}

func (s *service) sendAppdynamicsAlert(ctx context.Context, chatId uint32, message AppdynamicsMessage) (err error) {

	for _, event := range message.Events {

		if message.InvokeCallout {
			if "ERROR" == strings.ToUpper(event.Severity) {
				s.calloutService.InvokeCallout(ctx, chatId, "Appdynamics Critical Issue", event.EventMessage, true)
			}
		}

		message := event.EventMessage
		message = strings.Replace(message, "<b>", "*", -1)
		message = strings.Replace(message, "</b>", "*", -1)
		message = strings.Replace(message, "<br>", "\n", -1)

		var buffer bytes.Buffer
		switch strings.ToUpper(event.Severity) {
		case "ERROR":
			buffer.WriteString(emoji.Sprintf(":red_circle:"))
		case "WARN":
			buffer.WriteString(emoji.Sprintf(":warning:"))
		default:
			buffer.WriteString(emoji.Sprintf(":information_source:"))
		}

		buffer.WriteString(" ")
		buffer.WriteString(message)
		buffer.WriteString("\n")

		if event.Application.Name != "" {
			app := strings.Replace(event.Application.Name, "_", "\\_", -1)
			buffer.WriteString("*Application:* ")
			buffer.WriteString(app)
			buffer.WriteString("\n")
		}

		if event.Tier.Name != "" {
			ti := strings.Replace(event.Tier.Name, "_", "\\_", -1)
			buffer.WriteString("*Tier:* ")
			buffer.WriteString(ti)
			buffer.WriteString("\n")
		}

		if event.Node.Name != "" {
			no := strings.Replace(event.Node.Name, "_", "\\_", -1)
			buffer.WriteString("*Node:* ")
			buffer.WriteString(no)
			buffer.WriteString("\n")
		}

		log.Printf("Sending Alert %s", buffer.String())
		tmperr := s.alert.SendAlert(ctx, chatId, buffer.String())
		if tmperr != nil {
			err = tmperr
		}
	}
	return err
}

func (s *service) addAppdynamicsEndpoint(chat uint32, endpoint string) error {
	return s.store.addAppDynamicsEndpoint(chat, endpoint)
}

func (s *service) executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error {
	ipaddress, err := s.getIPAddressForNode(ctx, applicationID, nodeID, chatId)
	if err != nil {
		s.alert.SendError(ctx, err)
		return err
	}
	return s.ssh.ExecuteRemoteCommand(ctx, chatId, commandName, ipaddress)
}

func (s *service) getIPAddressForNode(ctx context.Context, application, node string, chat uint32) (string, error) {
	uri := fmt.Sprintf("/controller/rest/applications/%s/nodes/%s?output=json", application, node)
	response, err := doGet(ctx, uri, s.store, s.alert, chat)
	log.Println(response)
	if err != nil {
		s.alert.SendError(ctx, err)
		return "", err
	}

	var dat []interface{}
	err = json.Unmarshal([]byte(response), &dat)
	if err != nil {
		s.alert.SendError(ctx, err)
		return "", err
	}
	v := dat[0].(map[string]interface{})
	ipaddresses := v["ipAddresses"].(map[string]interface{})
	arrayIP := ipaddresses["ipAddresses"].([]interface{})
	for _, ipo := range arrayIP {
		ip := ipo.(string)
		if strings.Index(ip, ".") > 0 {
			return ip, nil
		}
	}
	return "", errors.New("no up address found")
}

// swagger:model
type AppdynamicsMessage struct {
	Environment string `json:"environment"`
	Policy      struct {
		TriggerTime string `json:"triggerTime"`
		Name        string `json:"name"`
	}
	Events        []Event `json:"events"`
	InvokeCallout bool    `json:"invoke_callout"`
}

type Event struct {
	Severity     string `json:"severity"`
	Application  Name
	Tier         Name
	Node         Name
	DisplayName  string `json:"displayName"`
	EventMessage string `json:"eventMessage"`
}

type Name struct {
	Name string `json:"name"`
}
