package telegram

import (
	"bytes"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
)

type Service interface {
	SendMessage(chatID int64, message string, messageID int) (err error)
	SendMessagePlainText(chatID int64, message string, messageID int) (err error)
	SendImageToGroup(image []byte, group int64) error
	SendKeyboard(buttons []string, text string, chat int64)
	RegisterCommand(command Command)
	RegisterCommandLet(commandlet Commandlet)
}

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	Execute(update tgbotapi.Update)
}

type Commandlet interface {
	CanExecute(update tgbotapi.Update, state State) bool
	Execute(update tgbotapi.Update, state State)
	NextState(update tgbotapi.Update, state State) string
	Fields(update tgbotapi.Update, state State) []string
}

type service struct {
	store Store
}

type commandCtor func() Command
type commandletCtor func() Commandlet

var commandList = map[string]commandCtor{}
var commandletList = []commandletCtor{}
var telegramBot *tgbotapi.BotAPI

func NewService(store Store) Service {
	s := &service{store}
	err := s.useBot(os.Getenv("BOT_KEY"))
	if err != nil {
		panic(err)
	}
	return s
}

func (s *service) SendMessage(chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID, message, messageID, true)
}

func (s *service) SendMessagePlainText(chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID, message, messageID, false)
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

func (s *service) RegisterCommand(command Command) {
	register(func() Command {
		return command
	})
}
func (s *service) RegisterCommandLet(commandlet Commandlet) {
	registerCommandlet(func() Commandlet {
		return commandlet
	})
}

func (s *service) SendKeyboard(buttons []string, text string, chat int64) {
	keyB := tgbotapi.NewReplyKeyboard()
	keyBRow := tgbotapi.NewKeyboardButtonRow()

	for i, l := range buttons {
		btn := tgbotapi.KeyboardButton{l, false, false}
		keyBRow = append(keyBRow, btn)
		if i > 0 && i%3 == 0 {
			keyB.Keyboard = append(keyB.Keyboard, keyBRow)
			keyBRow = tgbotapi.NewKeyboardButtonRow()
		}
	}
	if len(keyBRow) > 0 {
		keyB.Keyboard = append(keyB.Keyboard, keyBRow)
	}
	keyB.OneTimeKeyboard = true

	msg := tgbotapi.NewMessage(chat, text)
	msg.ReplyMarkup = keyB
	telegramBot.Send(msg)
}

func (s service) useBot(botkey string) error {
	var err error
	telegramBot, err = tgbotapi.NewBotAPI(botkey)
	if err != nil {
		return err
	}
	go func() {
		s.pollMessage()
	}()
	return nil

}

func (s service) pollMessage() {
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
				if s.executeCommandLet(update) {
					continue
				}
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			sendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID, false)
		}
	}
}

func sendMessage(chatID int64, message string, messageID int, markup bool) (err error) {

	log.Printf("Sending Message %s", message)

	msg := tgbotapi.NewMessage(chatID, message)
	if markup {
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
		go func() { command.Execute(update) }()
		return true
	}
	return false
}

func (s service) executeCommandLet(update tgbotapi.Update) bool {
	state := s.store.getState(update.Message.From.ID)
	for _, c := range commandletList {
		a := c()
		if a.CanExecute(update, state) {
			a.Execute(update, state)
			s.store.SetState(update.Message.From.ID, a.NextState(update, state), a.Fields(update, state))
			return true
		}
	}
	return false
}

func findCommand(command string) (a Command) {
	for _, item := range commandList {
		a = item()
		if a.CommandIdentifier() == command {
			return a
		}
	}
	return nil
}

func register(newfunc commandCtor) {
	id := newfunc().CommandIdentifier()
	commandList[id] = newfunc
}

func registerCommandlet(newFunc commandletCtor) {
	commandletList = append(commandletList, newFunc)
}

type help struct {
	telegram Service
}

func NewHelpCommand(telegram Service) Command {
	return &help{telegram}
}

func (s *help) CommandIdentifier() string {
	return "help"
}

func (s *help) CommandDescription() string {
	return "Gets list of commands"
}

func (s *help) Execute(update tgbotapi.Update) {
	var buffer bytes.Buffer
	for _, x := range getCommands() {
		buffer.WriteString("/")
		buffer.WriteString(x.Name)
		buffer.WriteString(" - ")
		buffer.WriteString(x.Description)
		buffer.WriteString("\n")
	}
	s.telegram.SendMessage(update.Message.Chat.ID, buffer.String(), update.Message.MessageID)
}

func getCommands() []commandDescription {
	result := make([]commandDescription, len(commandList))
	count := 0
	for _, x := range commandList {
		result[count] = commandDescription{x().CommandIdentifier(), x().CommandDescription()}
		count++
	}
	return result
}

/**
Basic information about a command
*/
type commandDescription struct {
	Name, Description string
}
