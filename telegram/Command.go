package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func init(){
	log.Println("Initialising command registry")
}

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	execute(update tgbotapi.Update)
}

type CommandDescription struct {
	Name, Description string
}

var commandList = []commandCtor{}

type commandCtor func() Command

func Register(newfund commandCtor) {
	commandList = append(commandList, newfund)
}

func findCommand(command string) (a Command) {
	for _, item := range commandList {
		a = item()
		if a.CommandIdentifier() == command {
			return a
		}
	};
	return nil
}

func ExecuteCommand(update tgbotapi.Update) {
	findCommand(update.Message.Command()).execute(update)
}

func GetCommands()  []CommandDescription{
	result := make([]CommandDescription,len(commandList))
	for i,x := range commandList {
		result[i] = CommandDescription{x().CommandIdentifier(), x().CommandDescription()}
	}
	return result
}
