package alert

import (
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/telegram"
)

type setGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setGroupCommand{service, store}
}

func (s *setGroupCommand) commandIdentifier() string {
	return "SetGroup"
}

func (s *setGroupCommand) commandDescription() string {
	return "Set Alert Group"
}

func (s *setGroupCommand) execute(update tgbotapi.Update) {
	s.store.setAlertGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(update.Message.Chat.ID, "group updated", update.Message.MessageID)
}

type setNonTechnicalGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetNonTechnicalGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setNonTechnicalGroupCommand{service, store}
}

func (s *setNonTechnicalGroupCommand) commandIdentifier() string {
	return "SetNonTechGroup"
}

func (s *setNonTechnicalGroupCommand) commandDescription() string {
	return "Set Non Technical Alert Group"
}

func (s *setNonTechnicalGroupCommand) execute(update tgbotapi.Update) {
	s.store.setNonTechnicalGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(update.Message.Chat.ID, "non technical group updated", update.Message.MessageID)
}

type setHeartbeatGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetHeartbeatGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setHeartbeatGroupCommand{service, store}
}

/* Set Heartbeat group */
func (s *setHeartbeatGroupCommand) commandIdentifier() string {
	return "SetHeartbeatGroup"
}

func (s *setHeartbeatGroupCommand) commandDescription() string {
	return "Set Heartbeat Group"
}

func (s *setHeartbeatGroupCommand) execute(update tgbotapi.Update) {
	s.store.setHeartbeatGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(update.Message.Chat.ID, "heartbeat group updated", update.Message.MessageID)
}
