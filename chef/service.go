package chef

import (
	json2 "encoding/json"
	"strings"
	"log"
	"bytes"
	"gopkg.in/kyokomi/emoji.v1"
	"fmt"
	"github.com/zamedic/go2hal/alert"
	"time"
	"github.com/go-chef/chef"
)

type Service interface {
	/*
	SendDeliveryAlert will unmarshal the input json and send an alert to the telegram group.
	*/
	SendDeliveryAlert(message string)
}

type service struct {
	alert alert.Service

	chefStore Store
}

func NewService() Service{

	return &service{}
}

/*
SendDeliveryAlert will unmarshal the input json and send an alert to the telegram group.
 */
func (s service)SendDeliveryAlert(message string) {
	var dat map[string]interface{}

	message = strings.Replace(message, "\n", "\\n", -1)

	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		s.alert.SendError(fmt.Errorf("delivery - error unmarshalling: %s", message))
		return
	}

	attachments := dat["attachments"].([]interface{})

	body := dat["text"].(string)
	bodies := strings.Split(body, "\n");
	url := bodies[0]
	url = strings.Replace(url, "<", "", -1)
	url = strings.Replace(url, ">", "", -1)

	parts := strings.Split(url, "|")

	//Loop though the attachmanets, there should be only 1
	var buffer bytes.Buffer
	buffer.WriteString(emoji.Sprint(":truck:"))
	buffer.WriteString(" ")
	buffer.WriteString("*chef Delivery*\n")

	if (strings.Contains(bodies[1], "failed")) {
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

	getfield(attachments, &buffer)

	buffer.WriteString("[")
	buffer.WriteString(parts[1])

	buffer.WriteString("](")
	buffer.WriteString(parts[0])
	buffer.WriteString(")")

	log.Printf("Sending Alert: %s", buffer.String())

	s.alert.SendAlert(buffer.String())

}
func getfield(attachments []interface{}, buffer *bytes.Buffer) {
	for _, attachment := range attachments {
		attachmentI := attachment.(map[string]interface{})
		fields := attachmentI["fields"].([]interface{})

		//Loop through the fields
		for _, field := range fields {
			fieldI := field.(map[string]interface{})
			buffer.WriteString("*")
			buffer.WriteString(fieldI["title"].(string))
			buffer.WriteString("* ")
			buffer.WriteString(fieldI["value"].(string))
			buffer.WriteString("\n")
		}
	}
}

func monitorQuarentined(s service) {
	for {
		checkQuarentined(s)
		time.Sleep(30 * time.Minute)
	}
}
func checkQuarentined(s service) {
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

	for _,r := range recipes {
		for _, e := range env {
			nodes := findNodesFromFriendlyNames(r.FriendlyName,e.FriendlyName,s)
			for _,n := range nodes {
				if strings.Index(n.environment,"quar") > 0 {
					s.alert.SendAlert(emoji.Sprintf(":hospital: *Node Quarantined* \n node %v has been placed in environment %v. Application %v ",n.name,strings.Replace(n.environment,"_", " ",-1), r.FriendlyName))
				}
			}
		}
	}

}

func findNodesFromFriendlyNames(recipe, environment string, s service) []node {
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

	client, err := getChefClient(s)
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

	result := make([]node, res.Total)

	for i, x := range res.Rows {
		s := x.(map[string]interface{})
		data := s["data"].(map[string]interface{})
		name := data["name"].(string)
		env := data["chef_environment"].(string)
		result[i] = node{name:name,environment:env}
	}

	return result

}

func getChefClient(s service) (client *chef.Client, err error ) {
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

type node struct {
	name        string
	environment string
}
