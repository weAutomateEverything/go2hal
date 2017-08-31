package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"github.com/zamedic/go2hal/database"
)

func init(){
	log.Println("Initializing Set Group Command")
	register(func() Command {
		return &setGroup{}
	})
}

type setGroup struct {

}

func (s *setGroup) CommandIdentifier() string {
	return "SetGroup"
}

func (s *setGroup) CommandDescription() string {
	return "Set Alert Group"
}

func (s *setGroup) execute(update tgbotapi.Update){
	database.SetAlertGroup(update.Message.Chat.ID)
	SendMessage(update.Message.Chat.ID,"group updated", update.Message.MessageID)
}


