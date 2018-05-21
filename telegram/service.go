package telegram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/weAutomateEverything/go2hal/auth"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Service interface {
	SendMessage(ctx context.Context, chatID int64, message string, messageID int) (err error)
	SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (err error)
	SendImageToGroup(ctx context.Context, image []byte, group int64) error
	SendKeyboard(ctx context.Context, buttons []string, text string, chat int64)
	RegisterCommand(command Command)
	RegisterCommandLet(commandlet Commandlet)
}

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	RestrictToAuthorised() bool
	Execute(update tgbotapi.Update)
}

type Commandlet interface {
	CanExecute(update tgbotapi.Update, state State) bool
	Execute(update tgbotapi.Update, state State)
	NextState(update tgbotapi.Update, state State) string
	Fields(update tgbotapi.Update, state State) []string
}

type service struct {
	store       Store
	authService auth.Service
}

type commandCtor func() Command
type commandletCtor func() Commandlet

var commandList = map[string]commandCtor{}
var commandletList = []commandletCtor{}
var telegramBot *tgbotapi.BotAPI

func NewService(store Store, authService auth.Service) Service {
	s := &service{store: store, authService: authService}
	key := os.Getenv("BOT_KEY")
	if key == "" {
		panic("BOT_KEY environment variable not set.")
	}
	err := s.useBot(key)
	if err != nil {
		panic(err.Error())
	}
	return s
}

func (s *service) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID, message, messageID, true)
}

func (s *service) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID, message, messageID, false)
}

func (s *service) SendImageToGroup(ctx context.Context, image []byte, group int64) error {
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

func (s *service) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) {
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

			if update.Message.NewChatMembers != nil {
				for _, user := range *update.Message.NewChatMembers {
					// Looks like the bot has been added to a new group - lets register the details.
					if user.ID == telegramBot.Self.ID {
						id, err := s.store.addBot(update.Message.Chat.ID)
						if err != nil {
							sendMessage(update.Message.Chat.ID, fmt.Sprintf("There was an error registering your bot: %v", err.Error()), update.Message.MessageID, false)
							continue
						}

						sendMessage(update.Message.Chat.ID, fmt.Sprintf("The bot has been successfully registered. Your token is %v", id), update.Message.MessageID, false)
						continue
					}
				}
			}

			if update.Message.IsCommand() {
				if s.executeCommand(update) {
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

func (s service) executeCommand(update tgbotapi.Update) bool {
	command := findCommand(update.Message.Command())
	if command != nil {
		if command.RestrictToAuthorised() {
			if !(s.authService.Authorize(strconv.Itoa(update.Message.From.ID))) {
				s.SendMessage(context.TODO(), update.Message.Chat.ID, "You are not authorized to use this transaction.", update.Message.MessageID)
				return false
			}
		}
		go func() { command.Execute(update) }()
		return true
	}
	return false
}

func (s service) executeCommandLet(update tgbotapi.Update) bool {
	state := s.store.getState(update.Message.From.ID, update.Message.Chat.ID)
	for _, c := range commandletList {
		a := c()
		if a.CanExecute(update, state) {
			a.Execute(update, state)
			s.store.SetState(update.Message.From.ID, update.Message.Chat.ID, a.NextState(update, state), a.Fields(update, state))
			return true
		}
	}
	return false
}

func findCommand(command string) (a Command) {
	for _, item := range commandList {
		a = item()
		if strings.ToLower(a.CommandIdentifier()) == strings.ToLower(command) {
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

func (s *help) RestrictToAuthorised() bool {
	return false
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
	s.telegram.SendMessage(context.TODO(), update.Message.Chat.ID, buffer.String(), update.Message.MessageID)
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
