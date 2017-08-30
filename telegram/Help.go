package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"bytes"
	"log"
)

func init() {
	log.Println("Initialising Help Command")
	Register(func() Command {
		return &Help{}
	})
}

type Help struct {
}

func (s *Help) CommandIdentifier() string {
	return "help"
}

func (s *Help) CommandDescription() string {
	return "Gets list of commands"
}

func (s *Help) execute(update tgbotapi.Update) {
	var buffer bytes.Buffer
	for _, x := range GetCommands(){
		buffer.WriteString(x.Name)
		buffer.WriteString(" - ")
		buffer.WriteString(x.Description)
		buffer.WriteString("\n")
	}
	SendMessage(update.Message.Chat.ID,buffer.String(),update.Message.MessageID)
}
