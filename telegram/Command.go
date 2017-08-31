package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func init(){
	log.Println("Initialising command registry")
}

type Command interface {
	commandIdentifier() string
	commandDescription() string
	execute(update tgbotapi.Update)
}

/**
Basic information about a command
 */
type CommandDescription struct {
	Name, Description string
}

var commandList = []commandCtor{}

type commandCtor func() Command

func register(newfund commandCtor) {
	commandList = append(commandList, newfund)
}

func findCommand(command string) (a Command) {
	for _, item := range commandList {
		a = item()
		if a.commandIdentifier() == command {
			return a
		}
	};
	return nil
}

func executeCommand(update tgbotapi.Update) {
	findCommand(update.Message.Command()).execute(update)
}

func getCommands()  []CommandDescription{
	result := make([]CommandDescription,len(commandList))
	for i,x := range commandList {
		result[i] = CommandDescription{x().commandIdentifier(), x().commandDescription()}
	}
	return result
}
