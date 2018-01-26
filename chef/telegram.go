package chef

import (
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/telegram"

	"github.com/go-chef/chef"
	"github.com/zamedic/go2hal/skynet"
)

/* Rebuild chef Node */
type rebuildChefNode struct {
	stateStore telegram.Store
	chefStore Store
	telegram telegram.Service
}

func (s *rebuildChefNode) commandIdentifier() string {
	return "RebuildChefNode"
}

func (s *rebuildChefNode) commandDescription() string {
	return "Rebuilds a node based on a chef search"

}

func (s *rebuildChefNode) execute(update tgbotapi.Update) {
	s.stateStore.SetState(update.Message.From.ID, "REBUILD_CHEF", nil)
	sendRecipeKeyboard(update.Message.Chat.ID, "Please select the application for the node you want to rebuild",s.telegram,s.chefStore)
}

/* Commandlets */

type rebuildChefNodeRecipeReply struct {
}

func (s *rebuildChefNodeRecipeReply) canExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "REBUILD_CHEF"
}

func (s *rebuildChefNodeRecipeReply) execute(update tgbotapi.Update, state telegram.State) {
	sendEnvironemtKeyboard(update.Message.Chat.ID, "Please select the environment of the node you want to rebuild")
}

func (s *rebuildChefNodeRecipeReply) nextState(update tgbotapi.Update, state telegram.State) string {
	return "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeRecipeReply) fields(update tgbotapi.Update, state telegram.State) []string {
	return []string{update.Message.Text}
}

/* ----------------------------- */

type rebuildChefNodeEnvironmentReply struct {
	telegram telegram.Service
}

func (s *rebuildChefNodeEnvironmentReply) canExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeEnvironmentReply) execute(update tgbotapi.Update, state telegram.State) {
	nodes := findNodesFromFriendlyNames(state.Field[0], update.Message.Text)
	res := make([]string,len(nodes))
	for i,x := range nodes {
		res[i] = x.name
	}
	s.telegram.SendKeyboard(res, "Select node to rebuild", update.Message.Chat.ID)
}

func (s *rebuildChefNodeEnvironmentReply) nextState(update tgbotapi.Update, state telegram.State) string {
	return "RebuildChefNodeSelectNode"
}

func (s *rebuildChefNodeEnvironmentReply) fields(update tgbotapi.Update, state telegram.State) []string {
	return append(state.Field, update.Message.Text)
}

/*------------------*/
type rebuildChefNodeExecute struct {
	skynet skynet.Skynet
}

func (s *rebuildChefNodeExecute) canExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "RebuildChefNodeSelectNode"
}

func (s *rebuildChefNodeExecute) execute(update tgbotapi.Update, state telegram.State) {
	go func() {
		err := RecreateNode(update.Message.Text, update.Message.From.FirstName)
		if err != nil {
			SendError(err)
		}
	}()
}

func (s *rebuildChefNodeExecute) nextState(update tgbotapi.Update, state telegram.State) string {
	return ""
}

func (s *rebuildChefNodeExecute) fields(update tgbotapi.Update, state database.State) []string {
	return nil
}

func sendRecipeKeyboard(chat int64, text string, telegram telegram.Service, chefStore chef.Store) {
	recipes, err := s.chefStore.GetRecipes()
	if err != nil {
		s.alert.SendError(err)
		return
	}

	l := make([]string, len(recipes))
	for x, i := range recipes {
		l[x] = i.FriendlyName
	}
	s.alert.sendKeyboard(l, text, chat)
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

