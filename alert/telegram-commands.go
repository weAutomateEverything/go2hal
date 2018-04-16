package alert

import (
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"gopkg.in/telegram-bot-api.v4"
)

type setGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setGroupCommand{service, store}
}

func (s *setGroupCommand) RestrictToAuthorised() bool {
	return true
}

func (s *setGroupCommand) CommandIdentifier() string {
	return "SetGroup"
}

func (s *setGroupCommand) CommandDescription() string {
	return "Set Alert Group"
}

func (s *setGroupCommand) Execute(update tgbotapi.Update) {
	s.store.setAlertGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(context.TODO(), update.Message.Chat.ID, "group updated", update.Message.MessageID)
}

type setNonTechnicalGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetNonTechnicalGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setNonTechnicalGroupCommand{service, store}
}

func (s *setNonTechnicalGroupCommand) RestrictToAuthorised() bool {
	return true
}

func (s *setNonTechnicalGroupCommand) CommandIdentifier() string {
	return "SetNonTechGroup"
}

func (s *setNonTechnicalGroupCommand) CommandDescription() string {
	return "Set Non Technical Alert Group"
}

func (s *setNonTechnicalGroupCommand) Execute(update tgbotapi.Update) {
	s.store.setNonTechnicalGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(context.TODO(), update.Message.Chat.ID, "non technical group updated", update.Message.MessageID)
}

type setHeartbeatGroupCommand struct {
	telegram telegram.Service
	store    Store
}

func NewSetHeartbeatGroupCommand(service telegram.Service, store Store) telegram.Command {
	return &setHeartbeatGroupCommand{service, store}
}

func (s *setHeartbeatGroupCommand) RestrictToAuthorised() bool {
	return true
}

/* Set Heartbeat group */
func (s *setHeartbeatGroupCommand) CommandIdentifier() string {
	return "SetHeartbeatGroup"
}

func (s *setHeartbeatGroupCommand) CommandDescription() string {
	return "Set Heartbeat Group"
}

func (s *setHeartbeatGroupCommand) Execute(update tgbotapi.Update) {
	s.store.setHeartbeatGroup(update.Message.Chat.ID)
	s.telegram.SendMessage(context.TODO(), update.Message.Chat.ID, "heartbeat group updated", update.Message.MessageID)
}
