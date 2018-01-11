package service

import (
	"gopkg.in/telegram-bot-api.v4"
	"bytes"
	"log"
)

func init() {
	log.Println("Initialising Help Command")
	register(func() command {
		return &help{}
	})
	log.Println("Initialising Help Command - Completed")

}

type help struct {
}

func (s *help) commandIdentifier() string {
	return "help"
}

func (s *help) commandDescription() string {
	return "Gets list of commands"
}

func (s *help) execute(update tgbotapi.Update) {
	var buffer bytes.Buffer
	for _, x := range getCommands(){
		buffer.WriteString(x.Name)
		buffer.WriteString(" - ")
		buffer.WriteString(x.Description)
		buffer.WriteString("\n")
	}
	SendMessage(update.Message.Chat.ID,buffer.String(),update.Message.MessageID)
}
