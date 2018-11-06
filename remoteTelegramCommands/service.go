package remoteTelegramCommands

import (
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
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

	s.telegram.RegisterCommand(newRemoteCommand(request.Name, request.Description, uint32(request.Group), response))
	for {
		time.Sleep(10 * time.Minute)
	}
}

type remoteCommand struct {
	name, help string
	grounp     uint32
	remote     RemoteCommand_RegisterCommandServer
}

func newRemoteCommand(name, help string, group uint32, remote RemoteCommand_RegisterCommandServer) telegram.Command {
	return &remoteCommand{name: name, help: help, remote: remote, grounp: group}
}

func (s remoteCommand) RestrictToAuthorised() bool {
	return true
}

func (s remoteCommand) CommandIdentifier() string {
	return s.name
}

func (s remoteCommand) CommandDescription() string {
	return s.help
}

func (s remoteCommand) GetCommandGroup() uint32 {
	return s.grounp
}

func (remoteCommand)Show(uint32) bool{
	return true
}

func (s remoteCommand) Execute(ctx context.Context, update tgbotapi.Update) {
	request := RemoteRequest{Message: update.Message.CommandArguments(), From: strconv.FormatInt(int64(update.Message.From.ID), 10)}
	s.remote.Send(&request)
}
