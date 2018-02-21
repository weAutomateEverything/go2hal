package chef

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"github.com/go-chef/chef"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/util"
	"gopkg.in/kyokomi/emoji.v1"
	"log"
	"strings"
	"time"
)

type Service interface {
	sendDeliveryAlert(message string)
	FindNodesFromFriendlyNames(recipe, environment string) []Node
}

type service struct {
	alert alert.Service

	chefStore Store
}

func NewService(alert alert.Service, chefStore Store) Service {
	s := &service{alert, chefStore}
	go func() {
		s.monitorQuarentined()
	}()
	return s
}

func (s *service) sendDeliveryAlert(message string) {
	var dat map[string]interface{}

	message = strings.Replace(message, "\n", "\\n", -1)

	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		s.alert.SendError(fmt.Errorf("delivery - error unmarshalling: %s", message))
		return
	}

	attachments := dat["attachments"].([]interface{})

	body := dat["text"].(string)
	bodies := strings.Split(body, "\n")
	url := bodies[0]
	url = strings.Replace(url, "<", "", -1)
	url = strings.Replace(url, ">", "", -1)

	parts := strings.Split(url, "|")

	//Loop though the attachmanets, there should be only 1
	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprint(":truck:"))
	buffer.WriteString(" ")
	buffer.WriteString("*chef Delivery*\n")

	if strings.Contains(bodies[1], "failed") {
		buffer.WriteString(emoji.Sprint(":interrobang:"))

	} else {
		switch bodies[1] {
		case "Delivered stage has completed for this change.":
			buffer.WriteString(emoji.Sprint(":+1:"))

		case "Change Delivered!":
			buffer.WriteString(emoji.Sprint(":white_check_mark:"))

		case "Acceptance Passed. Change is ready for delivery.":
			buffer.WriteString(emoji.Sprint(":ok_hand:"))

		case "Change Approved!":
			buffer.WriteString(emoji.Sprint(":white_check_mark:"))

		case "Verify Passed. Change is ready for review.":
			buffer.WriteString(emoji.Sprint(":mag_right:"))
		}
	}
	buffer.WriteString(" ")

	buffer.WriteString(bodies[1])
	buffer.WriteString("\n")

	util.Getfield(attachments, &buffer)

	buffer.WriteString("[")
	buffer.WriteString(parts[1])

	buffer.WriteString("](")
	buffer.WriteString(parts[0])
	buffer.WriteString(")")

	log.Printf("Sending Alert: %s", buffer.String())

	s.alert.SendAlert(buffer.String())

}

func (s *service) monitorQuarentined() {
	for {
		s.checkQuarentined()
		time.Sleep(30 * time.Minute)
	}
}
func (s *service) checkQuarentined() {
	recipes, err := s.chefStore.GetRecipes()
	if err != nil {
		s.alert.SendError(err)
		return
	}

	env, err := s.chefStore.GetChefEnvironments()
	if err != nil {
		s.alert.SendError(err)
		return
	}

	for _, r := range recipes {
		for _, e := range env {
			nodes := s.FindNodesFromFriendlyNames(r.FriendlyName, e.FriendlyName)
			for _, n := range nodes {
				if strings.Index(n.Environment, "quar") > 0 {
					s.alert.SendAlert(emoji.Sprintf(":hospital: *Node Quarantined* \n node %v has been placed in environment %v. Application %v ", n.Name, strings.Replace(n.Environment, "_", " ", -1), r.FriendlyName))
				}
			}
		}
	}

}

func (s *service) FindNodesFromFriendlyNames(recipe, environment string) []Node {
	chefRecipe, err := s.chefStore.GetRecipeFromFriendlyName(recipe)
	if err != nil {
		s.alert.SendError(err)
		return nil
	}

	chefEnv, err := s.chefStore.GetEnvironmentFromFriendlyName(environment)
	if err != nil {
		s.alert.SendError(err)
		return nil
	}

	client, err := s.getChefClient()
	if err != nil {
		s.alert.SendError(err)
		return nil
	}

	query, err := client.Search.NewQuery("node", fmt.Sprintf("recipe:%s AND chef_environment:%s", chefRecipe, chefEnv))
	if err != nil {
		s.alert.SendError(err)
		return nil
	}

	part := make(map[string]interface{})
	part["name"] = []string{"name"}
	part["chef_environment"] = []string{"chef_environment"}

	res, err := query.DoPartial(client, part)
	if err != nil {
		s.alert.SendError(err)
		return nil
	}

	result := make([]Node, res.Total)

	for i, x := range res.Rows {
		s := x.(map[string]interface{})
		data := s["data"].(map[string]interface{})
		name := data["name"].(string)
		env := data["chef_environment"].(string)
		result[i] = Node{Name: name, Environment: env}
	}

	return result

}

func (s *service) getChefClient() (client *chef.Client, err error) {
	c, err := s.chefStore.GetChefClientDetails()
	if err != nil {
		return nil, err
	}
	client, err = connect(c.Name, c.Key, c.URL)
	return client, err
}

func connect(name, key, url string) (client *chef.Client, err error) {
	client, err = chef.NewClient(&chef.Config{
		Name:    name,
		Key:     key,
		BaseURL: url,
		SkipSSL: true,
	})
	return
}

type Node struct {
	Name        string
	Environment string
}
