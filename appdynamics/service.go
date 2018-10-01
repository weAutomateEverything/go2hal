package appdynamics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/ssh"
	"gopkg.in/kyokomi/emoji.v1"
	"log"
	"strings"
)

type Service interface {
	sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) error
	addAppdynamicsEndpoint(chat uint32, endpoint string) error
	executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error
}

type service struct {
	alert alert.Service
	store Store
	ssh   ssh.Service
}

func NewService(alertService alert.Service, sshservice ssh.Service, store Store) Service {
	return &service{alert: alertService, ssh: sshservice, store: store}
}

func (s *service) sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) error {
	var m appdynamicsMessage

	message = strings.Replace(message, "\"\"", "\\\"\\\"", -1)

	err := json.Unmarshal([]byte(message), &m)

	if err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error unmarshalling App Dynamics Message: %s. error: %v", message, err))
		return err
	}

	for _, event := range m.Events {

		message := event.EventMessage
		message = strings.Replace(message, "<b>", "*", -1)
		message = strings.Replace(message, "</b>", "*", -1)
		message = strings.Replace(message, "<br>", "\n", -1)

		var buffer bytes.Buffer
		buffer.WriteString(emoji.Sprintf(":red_circle:"))
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
