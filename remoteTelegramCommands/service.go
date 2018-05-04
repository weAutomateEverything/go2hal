package remoteTelegramCommands

import (
	"github.com/weAutomateEverything/go2hal/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"time"
)

type service struct {
	telegram telegram.Service
}

func NewService(telegram telegram.Service) RemoteCommandServer {
	return &service{telegram: telegram}
}

func (s *service) RegisterCommand(request *RemoteCommandRequest, response RemoteCommand_RegisterCommandServer) error {
	s.telegram.RegisterCommand(newRemoteCommand(request.Name, request.Description, response))
	for {
		time.Sleep(10 * time.Minute)
	}
}

type remoteCommand struct {
	name, help string
	remote     RemoteCommand_RegisterCommandServer
}

func (s remoteCommand) RestrictToAuthorised() bool {
	return true
}

func newRemoteCommand(name, help string, remote RemoteCommand_RegisterCommandServer) telegram.Command {
	return &remoteCommand{name: name, help: help, remote: remote}

}

func (s remoteCommand) CommandIdentifier() string {
	return s.name
}

func (s remoteCommand) CommandDescription() string {
	return s.help
}

func (s remoteCommand) Execute(update tgbotapi.Update) {
	request := RemoteRequest{Message: update.Message.CommandArguments(), From: strconv.FormatInt(int64(update.Message.From.ID), 10)}
	s.remote.Send(&request)
}
func (s *service) RegisterCommandLet(request *Request, response RemoteCommand_RegisterCommandLetServer) error {
	s.telegram.RegisterCommandLet(searchChefRecipeReplyCommand(request.State, request.Nextstate, response))
	for {
		time.Sleep(10 * time.Minute)
	}
}
type searchChefRecipeReply struct {
	state, nextstate string
	remote           RemoteCommand_RegisterCommandLetServer
}

func searchChefRecipeReplyCommand(state, nextstate string, remote RemoteCommand_RegisterCommandLetServer) telegram.Commandlet {
	return &searchChefRecipeReply{state: state, nextstate: nextstate, remote: remote}

}

func (s *searchChefRecipeReply) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == s.state
}
func (s *searchChefRecipeReply) Execute(update tgbotapi.Update, state telegram.State) {
	request := Response{Message: update.Message.Text, From: strconv.FormatInt(int64(update.Message.From.ID), 10), Fields: state.Field}
	s.remote.Send(&request)

}

func (s *searchChefRecipeReply) NextState(update tgbotapi.Update, state telegram.State) string {
	return s.nextstate
}

func (s *searchChefRecipeReply) Fields(update tgbotapi.Update, state telegram.State) []string {
	return []string{update.Message.Text}
}