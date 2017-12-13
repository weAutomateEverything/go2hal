package service

import (
	"github.com/go-chef/chef"
	"github.com/zamedic/go2hal/database"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
)

func init(){
	chef, err :=database.IsChefConfigured()
	if err != nil {
		SendError(err)
	}
	if chef{
		register(func() command {
			return &rebuildChefNode{}
		})
		registerCommandlet(func() commandlet{
			return &rebuildChefNodeRecipeReply{}
		})
		registerCommandlet(func() commandlet{
			return &rebuildChefNodeEnvironmentReply{}
		})
		registerCommandlet(func() commandlet{
			return &rebuildChefNodeExecute{}
		})
	}
}

//AddChefClient Adds a chef client.
func AddChefClient(name, key, url string) error {
	//Check if the details work
	_, err := connect(name, key, url)
	if err != nil {
		return err
	}
	//No Error - therefore we assume a successful connection
	database.AddChefClient(name, url, key)

	return nil
}

func getChefClient() (client *chef.Client, err error) {
	c, err := database.GetChefClientDetails()
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

func findNodesNamesFromFriendlyNames(recipe, environment string) []string{
	chefRecipe, err := database.GetRecipeFromFriendlyName(recipe)
	if err != nil {
		SendError(err)
		return nil
	}

	chefEnv, err := database.GetEnvironmentFromFriendlyName(environment)
	if err != nil {
		SendError(err)
		return nil
	}

	client, err := getChefClient()
	if err != nil {
		SendError(err)
		return nil
	}

	query, err := client.Search.NewQuery("node", fmt.Sprintf("recipe:%s AND chef_environment:%s", chefRecipe, chefEnv))
	if err != nil {
		SendError(err)
		return nil
	}

	part := make(map[string]interface{})
	part["name"] = []string{"name"}

	res, err := query.DoPartial(client,part)
	if err != nil {
		SendError(err)
		return nil
	}

	result := make([]string, res.Total)

	for i, x := range res.Rows {
		s := x.(map[string]interface{})
		data := s["data"].(map[string]interface{})
		name := data["name"].(string)
		result[i] = name
	}

	return result

}

func sendRecipeKeyboard(chat int64, text string) {
	recipes, err := database.GetRecipes()
	if err != nil {
		SendError(err)
		return
	}

	l := make([]string, len(recipes))
	for x, i := range recipes {
		l[x] = i.FriendlyName
	}
	sendKeyboard(l, text, chat)
}

func sendEnvironemtKeyboard(chat int64, text string) {
	e, err := database.GetChefEnvironments()
	if err != nil {
		SendError(err)
		return
	}

	l := make([]string, len(e))
	for x, i := range e {
		l[x] = i.FriendlyName
	}
	sendKeyboard(l, text, chat)
}

/* Rebuild Chef Node */
type rebuildChefNode struct {
}

func (s *rebuildChefNode) commandIdentifier() string {
	return "RebuildChefNode"
}

func (s *rebuildChefNode) commandDescription() string {
	return "Rebuilds a node based on a chef search"

}

func (s *rebuildChefNode) execute(update tgbotapi.Update) {
	database.SetState(update.Message.From.ID, "REBUILD_CHEF", nil)
	sendRecipeKeyboard(update.Message.Chat.ID, "Please select the application for the node you want to rebuild")
}

/* Commandlets */

type rebuildChefNodeRecipeReply struct {
}

func (s *rebuildChefNodeRecipeReply) canExecute(update tgbotapi.Update, state database.State) bool {
	return state.State == "REBUILD_CHEF"
}

func (s *rebuildChefNodeRecipeReply) execute(update tgbotapi.Update, state database.State) {
	sendEnvironemtKeyboard(update.Message.Chat.ID, "Please select the environment of the node you want to rebuild")
}

func (s *rebuildChefNodeRecipeReply) nextState(update tgbotapi.Update, state database.State) string {
	return "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeRecipeReply) fields(update tgbotapi.Update, state database.State) []string {
	return []string{update.Message.Text}
}

/* ----------------------------- */

type rebuildChefNodeEnvironmentReply struct {
}

func (s *rebuildChefNodeEnvironmentReply) canExecute(update tgbotapi.Update, state database.State) bool {
	return state.State == "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeEnvironmentReply) execute(update tgbotapi.Update, state database.State) {
	nodes := findNodesNamesFromFriendlyNames(state.Field[0],update.Message.Text)
	sendKeyboard(nodes,"Select node to rebuild",update.Message.Chat.ID)
}

func (s *rebuildChefNodeEnvironmentReply) 	nextState(update tgbotapi.Update,state database.State) string{
	return "RebuildChefNodeSelectNode"
}


func (s *rebuildChefNodeEnvironmentReply) 	fields(update tgbotapi.Update,state database.State) []string {
	return append(state.Field,update.Message.Text)
}

/*------------------*/
type rebuildChefNodeExecute struct {

}

func (s *rebuildChefNodeExecute) canExecute(update tgbotapi.Update, state database.State) bool {
	return state.State == "RebuildChefNodeSelectNode"
}

func (s *rebuildChefNodeExecute) execute(update tgbotapi.Update, state database.State) {
	go func() {
		err := RecreateNode(update.Message.Text, update.Message.From.FirstName)
		if err != nil {
			SendError(err)
		}
	}()
}

func (s *rebuildChefNodeExecute) 	nextState(update tgbotapi.Update,state database.State) string{
	return ""
}


func (s *rebuildChefNodeExecute) 	fields(update tgbotapi.Update,state database.State) []string {
	return nil
}

