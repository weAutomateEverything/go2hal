package telegram

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/weAutomateEverything/go2hal/auth"
	"github.com/weAutomateEverything/halMessageClassification/api"
	"gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Service interface {
	SendMessage(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error)
	SendMessageWithCorrelation(ctx context.Context, chatID int64, message string, messageID int, correlationId string) (msgid int, err error)
	SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error)
	SendImageToGroup(ctx context.Context, image []byte, group int64) error
	SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) error

	SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (int, error)
	RegisterCommand(command Command)
	RegisterCommandLet(commandlet Commandlet)

	requestAuthorisation(ctx context.Context, chat uint32, name string) (string, error)
	pollAuthorisation(token string) (uint32, error)
}

type Command interface {
	CommandIdentifier() string
	CommandDescription() string
	RestrictToAuthorised() bool
	Execute(ctx context.Context, update tgbotapi.Update)
}

type Commandlet interface {
	CanExecute(update tgbotapi.Update, state State) bool
	Execute(ctx context.Context, update tgbotapi.Update, state State)
	NextState(update tgbotapi.Update, state State) string
	Fields(update tgbotapi.Update, state State) []string
}

type RemoteCommand interface {
	GetCommandGroup() uint32
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

func (s *service) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	return sendMessage(chatID, message, messageID, true)
}

func (s *service) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
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

func (s *service) SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) error {
	path := ""

	if runtime.GOOS == "windows" {
		path = "c:/temp/"
	} else {
		path = "/tmp/"
	}
	path = path + string(rand.Int()) + "." + extension
	err := ioutil.WriteFile(path, document, os.ModePerm)

	if err != nil {
		return err
	}

	msg := tgbotapi.NewDocumentUpload(group, path)

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

func (s *service) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (int, error) {
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
	m, err := telegramBot.Send(msg)
	if err != nil {
		return 0, err
	}
	return m.MessageID, nil
}

func (s service) requestAuthorisation(ctx context.Context, chat uint32, name string) (authtoken string, err error) {
	room, err := s.store.GetRoomKey(chat)
	if err != nil {
		return
	}
	id, err := s.SendKeyboard(ctx, []string{"Approve access", "Decline access"}, fmt.Sprintf("A user %v has requested access to edit the configuration for the room.", name), room)
	if err != nil {
		return "", err
	}

	token, err := s.store.newAuthRequest(id, room, name)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s service) pollAuthorisation(token string) (room uint32, err error) {
	roomId, err := s.store.useToken(token)
	if err != nil {
		return
	}

	return s.store.GetUUID(roomId, "")

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
						id, err := s.store.addBot(update.Message.Chat.ID, update.Message.Chat.Title)
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
				if update.Message.ReplyToMessage != nil {
					key, err := s.store.getCorrelationId(update.Message.Chat.ID, update.Message.ReplyToMessage.MessageID)
					if err != nil {
						log.Printf("Error getting correlation id %v", err)
						continue
					}
					err = s.store.SaveReply(update.Message.Chat.ID, update.Message.Text, key)
					if err != nil {
						log.Printf("Error saving reply %v", err)
						continue
					}
				}
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			sendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID, false)
		}
	}
}

func (s service) SendMessageWithCorrelation(ctx context.Context, chatID int64, message string, messageID int, correlationId string) (msgid int, err error) {
	msgid, err = s.SendMessage(ctx, chatID, message, messageID)
	if err != nil {
		return
	}

	err = s.store.saveMessageCorrelation(chatID, msgid, correlationId)
	return
}

func sendMessage(chatID int64, message string, messageID int, markup bool) (msgid int, err error) {

	log.Printf("Sending Message %s", message)

	msg := tgbotapi.NewMessage(chatID, message)
	if markup {
		msg.ParseMode = tgbotapi.ModeMarkdown
	}
	if messageID != 0 {
		msg.ReplyToMessageID = messageID
	}
	out, err := telegramBot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

	err = auditMessage(message, chatID, strconv.FormatInt(int64(out.MessageID), 10))
	if err != nil {
		log.Printf("Error auditing message: %v", err)
	}
	return out.MessageID, nil
}

func (s service) executeCommand(update tgbotapi.Update) bool {
	group, _ := s.store.GetUUID(update.Message.Chat.ID, update.Message.Chat.Title)
	command := findCommand(update.Message.Command(), group)
	if command != nil {
		if command.RestrictToAuthorised() {
			if !(s.authService.Authorize(strconv.Itoa(update.Message.From.ID))) {
				s.SendMessage(context.Background(), update.Message.Chat.ID, "You are not authorized to use this transaction.", update.Message.MessageID)
				return false
			}
		}
		go func() {
			command.Execute(context.Background(), update)
		}()
		return true
	}
	return false
}

func (s service) executeCommandLet(update tgbotapi.Update) bool {
	state := s.store.getState(update.Message.From.ID, update.Message.Chat.ID)
	for _, c := range commandletList {
		a := c()
		if a.CanExecute(update, state) {
			a.Execute(context.Background(), update, state)
			s.store.SetState(update.Message.From.ID, update.Message.Chat.ID, a.NextState(update, state), a.Fields(update, state))
			return true
		}
	}
	return false
}

func findCommand(command string, group uint32) (a Command) {
	//Check if its a standard command
	c, ok := commandList[strings.ToLower(command)]
	if ok {
		return c()
	}

	//see if its a command for a certain group
	c, ok = commandList[strings.ToLower(command)+strconv.FormatUint(uint64(group), 10)]

	if ok {
		return c()
	}

	return nil
}

func register(newfunc commandCtor) {
	id := newfunc().CommandIdentifier()
	remote, ok := newfunc().(RemoteCommand)
	if ok {
		id = id + strconv.FormatUint(uint64(remote.GetCommandGroup()), 10)
	}
	commandList[strings.ToLower(id)] = newfunc
}

func registerCommandlet(newFunc commandletCtor) {
	commandletList = append(commandletList, newFunc)
}

type help struct {
	telegram Service
	store    Store
}

func NewHelpCommand(telegram Service, store Store) Command {
	return &help{telegram, store}
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

func (s *help) Execute(ctx context.Context, update tgbotapi.Update) {
	var buffer bytes.Buffer
	for _, x := range getCommands() {
		if x.Group != 0 {
			g, err := s.store.GetRoomKey(x.Group)
			if err != nil {
				continue
			}
			if g != update.Message.Chat.ID {
				continue
			}
		}
		buffer.WriteString("/")
		buffer.WriteString(x.Name)
		buffer.WriteString(" - ")
		buffer.WriteString(x.Description)
		buffer.WriteString("\n")
	}
	s.telegram.SendMessage(ctx, update.Message.Chat.ID, buffer.String(), update.Message.MessageID)
}

func getCommands() []commandDescription {
	result := make([]commandDescription, len(commandList))
	count := 0
	for _, x := range commandList {
		r, ok := x().(RemoteCommand)
		if ok {
			result[count] = commandDescription{x().CommandIdentifier(), x().CommandDescription(), r.GetCommandGroup()}
		} else {
			result[count] = commandDescription{x().CommandIdentifier(), x().CommandDescription(), 0}
		}
		count++
	}
	return result
}

func auditMessage(message string, chat int64, messageId string) error {
	if strings.Contains(message, "Negative Sentiment Message Detected for group") {
		log.Println("Ignoring message")
		return nil
	}
	log.Printf("Checking Audit endpoint: %v", os.Getenv("HAL_API_SERVICES"))
	if os.Getenv("HAL_API_SERVICES") != "" {
		req := api.TextEvent{
			Message:   message,
			MessageID: messageId,
			Chat:      chat,
		}

		b, err := json.Marshal(req)
		if err != nil {
			return err
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		resp, err := http.Post(fmt.Sprintf("%v/halMessageClassification", os.Getenv("HAL_API_SERVICES")), "application/json", bytes.NewReader(b))
		if err != nil {
			return err
		}

		resp.Body.Close()
	}
	return nil
}

/**
Basic information about a command
*/
type commandDescription struct {
	Name, Description string
	Group             uint32
}
