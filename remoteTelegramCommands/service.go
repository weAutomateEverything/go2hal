package remoteTelegramCommands

import (
	"github.com/zamedic/go2hal/telegram"
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
	request := RemoteRequest{Message: update.Message.Text, From: strconv.FormatInt(int64(update.Message.From.ID), 10)}
	s.remote.Send(&request)
}
