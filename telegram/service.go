package telegram

import (
	"runtime"
	"io/ioutil"
	"os"
	"gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"log"
	"time"
	"bytes"
)

type Service interface {
	SendMessage(chatID int64, message string, messageID int) (err error)
	SendImageToGroup(image []byte, group int64) error
	SendKeyboard(buttons []string, text string, chat int64)
	RegisterCommand(command Command)
	RegisterCommandLet(commandlet Commandlet)

}

type Command interface {
	commandIdentifier() string
	commandDescription() string
	execute(update tgbotapi.Update)
}

type Commandlet interface {
	canExecute(update tgbotapi.Update, state State) bool
	execute(update tgbotapi.Update, state State)
	nextState(update tgbotapi.Update, state State) string
	fields(update tgbotapi.Update, state State) []string
}

type service struct {
	store Store
}

type commandCtor func() Command
type commandletCtor func() Commandlet

var commandList = []commandCtor{}
var commandletList = []commandletCtor{}
var telegramBot *tgbotapi.BotAPI

func NewService(store Store) Service{
	findFreeBot(store)
	return &service{store}
}


func (s *service) SendMessage(chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID, message, messageID, true)
}

func (s *service) SendImageToGroup(image []byte, group int64) error {
	path := ""

	if runtime.GOOS == "windows" {
		path = "c:/temp/"
	} else {
		path = "/tmp/"
	}
	path = path + string(rand.Int()) + ".png"
	err := ioutil.WriteFile(path, image, os.ModePerm)

	if err != nil {
		return err
	}

	msg := tgbotapi.NewPhotoUpload(group, path)

	_, err = telegramBot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *service)RegisterCommand(command Command){
	register(func() Command{
		return &command
	})
}
func (s *service)RegisterCommandLet(commandlet  Commandlet){
	registerCommandlet(func() Commandlet{
		return &commandlet
	})
}

func (s *service) SendKeyboard(buttons []string, text string, chat int64){
	keyB := tgbotapi.NewReplyKeyboard()
	keyBRow := tgbotapi.NewKeyboardButtonRow()

	for i,l := range buttons {
		btn := tgbotapi.KeyboardButton{l,false,false}
		keyBRow = append(keyBRow,btn)
		if i > 0 && i % 3 == 0 {
			keyB.Keyboard = append(keyB.Keyboard, keyBRow)
			keyBRow = tgbotapi.NewKeyboardButtonRow()
		}
	}
	if len(keyBRow) > 0 {
		keyB.Keyboard = append(keyB.Keyboard, keyBRow)
	}
	keyB.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(chat,text)
	msg.ReplyMarkup = keyB
	telegramBot.Send(msg)
}

func findFreeBot(s Store) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	for true {
		bots := s.listBots()
		if bots != nil {
			for _, botkey := range bots {
				useBot(botkey,s)
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func useBot(botkey string, s Store) error{
	var err error
	telegramBot, err = tgbotapi.NewBotAPI(botkey)
	if err != nil {
		return err
	}
	go func() {
		pollMessage(s)
	}()
	return nil

}

func pollMessage(s Store){
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for true {
		updates, err := telegramBot.GetUpdates(u)
		if err != nil {
			log.Panic(err)
		}

		for _, update := range updates {
			if update.UpdateID >= u.Offset {
				u.Offset = update.UpdateID + 1
			}

			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() {
				if executeCommand(update) {
					continue
				}
			}

			if update.Message != nil {
				if executeCommandLet(update,s){
					continue
				}
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			sendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID,false)
		}
	}
}

func sendMessage(chatID int64, message string, messageID int, markup bool) (err error) {

	log.Printf("Sending Message %s", message)

	msg := tgbotapi.NewMessage(chatID, message)
	if (markup) {
		msg.ParseMode = tgbotapi.ModeMarkdown
	}
	if messageID != 0 {
		msg.ReplyToMessageID = messageID
	}
	_, err = telegramBot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func executeCommand(update tgbotapi.Update) bool {
	command := findCommand(update.Message.Command())
	if command != nil {
		go func() {command.execute(update)}()
		return true
	}
	return false
}

func executeCommandLet(update tgbotapi.Update, s Store) bool {
	state := s.getState(update.Message.From.ID)
	for _, c := range commandletList{
		a := c()
		if a.canExecute(update,state) {
			a.execute(update,state)
			s.SetState(update.Message.From.ID,a.nextState(update,state),a.fields(update,state))
			return true
		}
	}
	return false
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

func register(newfunc commandCtor) {
	commandList = append(commandList, newfunc)
}

func registerCommandlet(newFunc commandletCtor){
	commandletList = append(commandletList,newFunc)
}

type help struct {
	telegram Service
}

func NewHelpCommand(telegram Service) Command {
	return &help{telegram}
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
		buffer.WriteString("/")
		buffer.WriteString(x.Name)
		buffer.WriteString(" - ")
		buffer.WriteString(x.Description)
		buffer.WriteString("\n")
	}
	s.telegram.SendMessage(update.Message.Chat.ID,buffer.String(),update.Message.MessageID)
}



func getCommands() []commandDescription {
	result := make([]commandDescription, len(commandList))
	for i, x := range commandList {
		result[i] = commandDescription{x().commandIdentifier(), x().commandDescription()}
	}
	return result
}

/**
Basic information about a command
 */
type commandDescription struct {
	Name, Description string
}



