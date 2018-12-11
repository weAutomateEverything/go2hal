package ssh

import (
	"context"
	"github.com/weAutomateEverything/go2hal/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type cmd struct {
	store      Store
	stateStore telegram.Store
	telegram   telegram.Service
}

func NewCommand(store Store,
	stateStore telegram.Store,
	telegram telegram.Service) telegram.Command {
	return &cmd{
		store:      store,
		telegram:   telegram,
		stateStore: stateStore,
	}
}

func (cmd) CommandIdentifier() string {
	return "execute"
}

func (cmd) CommandDescription() string {
	return "run a predefined ssh command"
}

func (cmd) RestrictToAuthorised() bool {
	return true
}

func (s cmd) Show(chat uint32) bool {
	key, err := s.store.getKey(chat)
	if err != nil {
		log.Println(err)
		return false
	}
	if key == nil {
		return false
	}

	commands, err := s.store.getCommands(chat)
	if err != nil {
		log.Println(err)
		return false
	}
	if commands == nil {
		return false
	}

	servers, err := s.store.getServers(chat)
	if err != nil {
		log.Println(err)
		return false
	}
	if servers == nil {
		return false
	}

	return true
}

func (s cmd) Execute(ctx context.Context, update tgbotapi.Update) {
	id, err := s.stateStore.GetUUID(update.Message.Chat.ID, update.Message.Chat.Description)
	if err != nil {
		log.Println(err)
		s.telegram.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

	commands, err := s.store.getCommands(id)
	if err != nil {
		log.Println(err)
		s.telegram.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

	keys := make([]string, len(commands))
	for x, command := range commands {
		keys[x] = command.Name
	}
	_, err = s.telegram.SendKeyboard(ctx, keys, "Which command would you like to run? ", update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		s.telegram.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

	err = s.stateStore.SetState(update.Message.From.ID, update.Message.Chat.ID, "SSH_COMMAND", nil)
	if err != nil {
		log.Println(err)
		s.telegram.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

}
