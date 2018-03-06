package skynet

import (
	"errors"
	"fmt"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/chef"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v4"
	"runtime/debug"
)

type rebuildNode struct {
	alertService  alert.Service
	skynetService Service
}

func NewRebuildNodeCommand(alertService alert.Service, skynetService Service) telegram.Command {
	return &rebuildNode{alertService, skynetService}
}

/* Rebuild Node */
func (s *rebuildNode) CommandIdentifier() string {
	return "RebuildNode"
}

func (s *rebuildNode) CommandDescription() string {
	return "Rebuilds a node"
}

func (s *rebuildNode) Execute(update tgbotapi.Update) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			s.alertService.SendError(context.TODO(), errors.New(fmt.Sprint(err)))
			s.alertService.SendError(context.TODO(), errors.New(string(debug.Stack())))

		}
	}()
	s.skynetService.RecreateNode(context.TODO(), update.Message.CommandArguments(), update.Message.From.UserName)
}

/* ------------------- */

/* Rebuild chef Node */
type rebuildChefNode struct {
	stateStore telegram.Store
	chefStore  chef.Store
	alert      alert.Service
	telegram   telegram.Service
}

func NewRebuildCHefNodeCommand(stateStore telegram.Store, chefStore chef.Store, telegram telegram.Service,
	alert alert.Service) telegram.Command {
	return &rebuildChefNode{stateStore, chefStore, alert, telegram}
}

func (s *rebuildChefNode) CommandIdentifier() string {
	return "RebuildChefNode"
}

func (s *rebuildChefNode) CommandDescription() string {
	return "Rebuilds a node based on a chef search"

}

func (s *rebuildChefNode) Execute(update tgbotapi.Update) {
	s.stateStore.SetState(update.Message.From.ID, "REBUILD_CHEF", nil)
	sendRecipeKeyboard(update.Message.Chat.ID, "Please select the application for the node you want to rebuild",
		s.alert, s.chefStore, s.telegram)
}

/* Commandlets */

type rebuildChefNodeRecipeReply struct {
	store    chef.Store
	alert    alert.Service
	telegram telegram.Service
}

func NewRebuildChefNodeRecipeReplyCommandlet(store chef.Store, alert alert.Service,
	telegram telegram.Service) telegram.Commandlet {
	return &rebuildChefNodeRecipeReply{store, alert, telegram}
}

func (s *rebuildChefNodeRecipeReply) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "REBUILD_CHEF"
}

func (s *rebuildChefNodeRecipeReply) Execute(update tgbotapi.Update, state telegram.State) {
	sendEnvironemtKeyboard(update.Message.Chat.ID, "Please select the environment of the node you want to rebuild", s.store, s.alert, s.telegram)
}

func (s *rebuildChefNodeRecipeReply) NextState(update tgbotapi.Update, state telegram.State) string {
	return "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeRecipeReply) Fields(update tgbotapi.Update, state telegram.State) []string {
	return []string{update.Message.Text}
}

/* ----------------------------- */

type rebuildChefNodeEnvironmentReply struct {
	telegram    telegram.Service
	service     Service
	chefService chef.Service
}

func NewRebuildChefNodeEnvironmentReplyCommandlet(telegram telegram.Service, service Service,
	chefService chef.Service) telegram.Commandlet {
	return &rebuildChefNodeEnvironmentReply{telegram, service, chefService}
}

func (s *rebuildChefNodeEnvironmentReply) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "RebuildChefNodeEnvironment"
}

func (s *rebuildChefNodeEnvironmentReply) Execute(update tgbotapi.Update, state telegram.State) {
	nodes := s.chefService.FindNodesFromFriendlyNames(state.Field[0], update.Message.Text)
	res := make([]string, len(nodes))
	for i, x := range nodes {
		res[i] = x.Name
	}
	s.telegram.SendKeyboard(context.TODO(), res, "Select node to rebuild", update.Message.Chat.ID)
}

func (s *rebuildChefNodeEnvironmentReply) NextState(update tgbotapi.Update, state telegram.State) string {
	return "RebuildChefNodeSelectNode"
}

func (s *rebuildChefNodeEnvironmentReply) Fields(update tgbotapi.Update, state telegram.State) []string {
	return append(state.Field, update.Message.Text)
}

/*------------------*/
type rebuildChefNodeExecute struct {
	skynet Service
	alert  alert.Service
}

func NewRebuildChefNodeExecute(skynet Service, alert alert.Service) telegram.Commandlet {
	return &rebuildChefNodeExecute{skynet, alert}
}

func (s *rebuildChefNodeExecute) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "RebuildChefNodeSelectNode"
}

func (s *rebuildChefNodeExecute) Execute(update tgbotapi.Update, state telegram.State) {
	go func() {
		err := s.skynet.RecreateNode(context.TODO(), update.Message.Text, update.Message.From.FirstName)
		if err != nil {
			s.alert.SendError(context.TODO(), err)
		}
	}()
}

func (s *rebuildChefNodeExecute) NextState(update tgbotapi.Update, state telegram.State) string {
	return ""
}

func (s *rebuildChefNodeExecute) Fields(update tgbotapi.Update, state telegram.State) []string {
	return nil
}

func sendRecipeKeyboard(chat int64, text string, alert alert.Service, chefStore chef.Store, telegram telegram.Service) {
	recipes, err := chefStore.GetRecipes()
	if err != nil {
		alert.SendError(context.TODO(), err)
		return
	}

	l := make([]string, len(recipes))
	for x, i := range recipes {
		l[x] = i.FriendlyName
	}
	telegram.SendKeyboard(context.TODO(), l, text, chat)
}

func sendEnvironemtKeyboard(chat int64, text string, store chef.Store, alert alert.Service, telegram telegram.Service) {
	e, err := store.GetChefEnvironments()
	if err != nil {
		alert.SendError(context.TODO(), err)
		return
	}

	l := make([]string, len(e))
	for x, i := range e {
		l[x] = i.FriendlyName
	}
	telegram.SendKeyboard(context.TODO(), l, text, chat)
}
