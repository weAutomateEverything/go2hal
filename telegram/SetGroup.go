package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func init(){
	log.Println("Initializing Set Group Command")
	Register(func() Command {
		return &SetGroup{}
	})
}

type SetGroup struct {

}

func (s *SetGroup) CommandIdentifier() string {
	return "SetGroup"
}

func (s *SetGroup) CommandDescription() string {
	return "Set Alert Group"
}

func (s *SetGroup) execute(update tgbotapi.Update){

}


