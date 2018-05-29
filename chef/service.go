package chef

import (
	"bytes"
	"context"
	"encoding/base64"
	json2 "encoding/json"
	"fmt"
	"github.com/go-chef/chef"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/util"
	"gopkg.in/kyokomi/emoji.v1"
	"os"
	"strings"
	"time"
)

type Service interface {
	sendDeliveryAlert(ctx context.Context, chatId uint32, message string)
	FindNodesFromFriendlyNames(recipe, environment string, chat uint32) []Node

	getAllRecipes() ([]string, error)
	getRecipesForGroup(group uint32) ([]Recipe, error)
	addRecipeToGroup(ctx context.Context, group uint32, recipeName, friendlyName string) error

	getEnvironmentForGroup(group uint32) ([]ChefEnvironment, error)
	addEnvironmentToGroup(group uint32, name string, friendlyname string) error
}

func NewService(alert alert.Service, chefStore Store) Service {
	s := &service{alert, chefStore}
	go func() {
		s.monitorQuarentined()
	}()
	return s
}

type service struct {
	alert alert.Service

	chefStore Store
}

func (s *service) getEnvironmentForGroup(group uint32) ([]ChefEnvironment, error) {
	env, err := s.chefStore.GetEnvironmentForGroup(group)
	if err != nil {
		return nil, err
	}
	if env == nil {
		env = make([]ChefEnvironment, 0)
	}
	return env, err
}

func (s *service) addEnvironmentToGroup(group uint32, name string, friendlyname string) error {
	err := s.chefStore.AddChefEnvironment(name, friendlyname, group)
	if err != nil {
		return err
	}
	s.alert.SendAlert(context.TODO(), group, emoji.Sprintf(":new: Chef Environment %v has been added to the group for monitoring. The friendly name for the recipe is %v", name, friendlyname))
	return nil
}

func (s *service) addRecipeToGroup(ctx context.Context, group uint32, recipeName, friendlyName string) error {
	err := s.chefStore.AddRecipe(recipeName, friendlyName, group)
	if err != nil {
		return err
	}
	s.alert.SendAlert(ctx, group, emoji.Sprintf(":new: Chef Recipe %v has been added to the group for monitoring. The friendly name for the recipe is %v", recipeName, friendlyName))
	return nil
}

func (s *service) getAllRecipes() (result []string, err error) {
	c, err := s.getChefClient()
	if err != nil {
		return
	}
	query, err := c.Search.NewQuery("cookbooks", "name:*")
	if err != nil {
		return
	}

	part := make(map[string]interface{})
	part["name"] = []string{"name"}

	res, err := query.DoPartial(c, part)

	if err != nil {
		return
	}

	result = make([]string, res.Total)

	for i, x := range res.Rows {
		s := x.(map[string]interface{})
		data := s["data"].(map[string]interface{})
		name := data["name"].(string)
		result[i] = name
	}

	return

}

func (s *service) getRecipesForGroup(group uint32) (result []Recipe, err error) {
	result, err = s.chefStore.GetRecipesForGroup(group)
	if result == nil && err != nil {
		result = make([]Recipe, 0)
	}
	return
}

func (s *service) sendDeliveryAlert(ctx context.Context, chatId uint32, message string) {
	var dat map[string]interface{}

	message = strings.Replace(message, "\n", "\\n", -1)

	if err := json2.Unmarshal([]byte(message), &dat); err != nil {
		s.alert.SendError(ctx, fmt.Errorf("delivery - error unmarshalling: %s", message))
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

	if len(bodies) > 1 {
		buildDeliveryEnent(&buffer, bodies[1])
	} else {
		buffer.WriteString(emoji.Sprintf(":rage1: New Code Review \n"))
	}

	util.Getfield(attachments, &buffer)

	buffer.WriteString("[")
	buffer.WriteString(parts[1])

	buffer.WriteString("](")
	buffer.WriteString(parts[0])
	buffer.WriteString(")")

	s.alert.SendAlert(ctx, chatId, buffer.String())

}

func buildDeliveryEnent(buffer *bytes.Buffer, body string) {
	if strings.Contains(body, "failed") {
		buffer.WriteString(emoji.Sprint(":interrobang:"))

	} else {
		switch body {
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

	buffer.WriteString(body)
	buffer.WriteString("\n")
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
		s.alert.SendError(context.TODO(), err)
		return
	}

	for _, r := range recipes {
		env, err := s.chefStore.GetEnvironmentForGroup(r.ChatID)
		if err != nil {
			s.alert.SendError(context.TODO(), err)
			continue
		}
		for _, e := range env {
			nodes := s.FindNodesFromFriendlyNames(r.FriendlyName, e.FriendlyName, r.ChatID)
			for _, n := range nodes {
				if strings.Index(n.Environment, "quar") > 0 {
					//We have found a quarentined Node - Now we need to check the recipes and environment to find out who wants to know about this, in this environment
					s.alert.SendAlert(context.TODO(), r.ChatID, emoji.Sprintf(":hospital: *Node Quarantined* \n node %v has been placed in environment %v. Application %v ", n.Name, strings.Replace(n.Environment, "_", " ", -1), r.FriendlyName))
				}
			}
		}
	}

}

func (s *service) FindNodesFromFriendlyNames(recipe, environment string, chat uint32) []Node {
	chefRecipe, err := s.chefStore.GetRecipeFromFriendlyName(recipe, chat)
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return nil
	}

	chefEnv, err := s.chefStore.GetEnvironmentFromFriendlyName(environment, chat)
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return nil
	}

	client, err := s.getChefClient()
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return nil
	}

	query, err := client.Search.NewQuery("node", fmt.Sprintf("recipe:%s AND chef_environment:%s", chefRecipe, chefEnv))
	if err != nil {
		s.alert.SendError(context.TODO(), err)
		return nil
	}

	part := make(map[string]interface{})
	part["name"] = []string{"name"}
	part["chef_environment"] = []string{"chef_environment"}

	res, err := query.DoPartial(client, part)
	if err != nil {
		s.alert.SendError(context.TODO(), err)
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
	client, err = connect(getChefUserName(), getChefUserKey(), getChefURL())
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

func getChefUserName() string {
	return os.Getenv("CHEF_USER")
}

func getChefUserKey() string {
	s := os.Getenv("CHEF_KEY")
	r, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return string(r)
}

func getChefURL() string {
	return os.Getenv("CHEF_URL")
}
