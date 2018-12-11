package ssh

import (
	"context"
	"github.com/weAutomateEverything/go2hal/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type sshSelect struct {
	store           Store
	telegramStore   telegram.Store
	telegramService telegram.Service
}

func NewSshSelectCommand(store Store,
	telegramStore telegram.Store,
	telegramService telegram.Service) telegram.Commandlet {
	return &sshSelect{
		telegramService: telegramService,
		store:           store,
		telegramStore:   telegramStore,
	}
}

func (sshSelect) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "SSH_COMMAND"
}

func (s sshSelect) Execute(ctx context.Context, update tgbotapi.Update, state telegram.State) {
	id, err := s.telegramStore.GetUUID(update.Message.Chat.ID, update.Message.Chat.Description)
	if err != nil {
		log.Println(err)
		s.telegramService.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

	servers, err := s.store.getServers(id)
	if err != nil {
		log.Println(err)
		s.telegramService.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}

	buttons := make([]string, len(servers))

	for x, server := range servers {
		buttons[x] = server.Description
	}

	_, err = s.telegramService.SendKeyboard(ctx, buttons, "On which server should this command run? ", update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		s.telegramService.SendMessage(ctx, update.Message.Chat.ID, "Technical Error - please let the devs know", update.Message.MessageID)
		return
	}
}

func (sshSelect) NextState(update tgbotapi.Update, state telegram.State) string {
	return "SSH_SERVER"
}

func (sshSelect) Fields(update tgbotapi.Update, state telegram.State) []string {
	return []string{update.Message.Text}
}
