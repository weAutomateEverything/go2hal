package ssh

import (
	"context"
	"github.com/weAutomateEverything/go2hal/telegram"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type sshExecute struct {
	service         Service
	store           Store
	telegramStore   telegram.Store
	telegramService telegram.Service
}

func NewSSHExecute(service Service,
	store Store,
	telegramStore telegram.Store,
	telegramService telegram.Service) telegram.Commandlet {
	return &sshExecute{
		telegramStore:   telegramStore,
		store:           store,
		service:         service,
		telegramService: telegramService,
	}
}

func (sshExecute) CanExecute(update tgbotapi.Update, state telegram.State) bool {
	return state.State == "SSH_SERVER"
}

func (s sshExecute) Execute(ctx context.Context, update tgbotapi.Update, state telegram.State) {
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

	for _, server := range servers {
		if server.Description == update.Message.Text {
			err = s.service.ExecuteRemoteCommand(ctx, id, state.Field[0], server.Address)

		}
	}

}

func (sshExecute) NextState(update tgbotapi.Update, state telegram.State) string {
	return ""
}

func (sshExecute) Fields(update tgbotapi.Update, state telegram.State) []string {
	return nil
}
