package service

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/database"
	"time"
	"gopkg.in/kyokomi/emoji.v1"
	"fmt"
)

//HalBot Structure to describe the state of the bot
type HalBot struct {
	Running bool
	bot     *tgbotapi.BotAPI
}

type command interface {
	commandIdentifier() string
	commandDescription() string
	execute(update tgbotapi.Update)
}

/**
Basic information about a command
 */
type commandDescription struct {
	Name, Description string
}
type commandCtor func() command

type setHeartbeatGroup struct {
}

var commandList = []commandCtor{}
var hal *HalBot;

var bot *tgbotapi.BotAPI
var err error

func init() {
	hal = &HalBot{false, nil}
	go func() {
		findFreeBot()
	}()
	go func() {
		pollForMessages()
	}()
	go func() {
		heartbeat()
	}()
	register(func() command {
		return &setHeartbeatGroup{}
	})
}

/*
GetBot returns a handler back to the bot
 */
func GetBot() *HalBot {
	return hal
}

/*
SendMessage sends a test message to the chat id.
 */
func SendMessage(chatID int64, message string, messageID int) (err error) {
	return sendMessage(chatID,message,messageID,true)
}


func sendMessage(chatID int64, message string, messageID int, markup bool) (err error) {
	if !hal.Running {
		log.Println("Unable to send message as no bot is connected. ")
		return nil
	}
	log.Printf("Sending Message %s", message)
	if (!hal.Running) {
		database.AddMessageToQueue(message, chatID, messageID)
	}
	msg := tgbotapi.NewMessage(chatID, message)
	if (markup) {
		msg.ParseMode = tgbotapi.ModeMarkdown
	}
	if messageID != 0 {
		msg.ReplyToMessageID = messageID
	}
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
	} else {
		database.SendMessage()
	}
	return nil
}

func SendError(err error) {
	log.Println(err.Error())
	sendToHeartbeatGroup(emoji.Sprintf(":poop: %s %s", bot.Self.UserName, err.Error()))
}

func sendToHeartbeatGroup(message string) {
	chatId, err := database.HeartbeatGroup()
	if err == nil && chatId != 0 {
		sendMessage(chatId, message, 0,false)
	} else {
		log.Printf("Could not send %s to heartbeat group", message)
	}
}

/*
TestBot checks if the token is for a valid bot.
 */
func TestBot(token string) error {
	_, err = tgbotapi.NewBotAPI(token)
	return err
}

func findFreeBot() {
	for true {
		bots := database.ListBots()
		if bots != nil {
			for _, botkey := range bots {
				useBot(botkey)
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func useBot(botkey string) {
	bot, err = tgbotapi.NewBotAPI(botkey)
	if err != nil {
		hal.Running = false
		log.Printf("Error getting bot token: %s", err.Error())
		return
	}
	hal.Running = true
	hal.bot = bot

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	sendToHeartbeatGroup(fmt.Sprintf("%s back online", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for true {
		updates, err := bot.GetUpdates(u)
		if err != nil {
			log.Println("Releasing bot ", bot.Self.UserName)
			log.Println(err.Error())
			hal.Running = false
			hal.bot = nil
			return
		}
		database.HeartbeatBot(botkey, bot.Self.UserName)

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
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			SendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
		}
	}
}

func pollForMessages() {
	for true {
		if (hal.Running) {
			messages := database.GetMessages()
			for _, x := range messages {
				SendMessage(x.ChatID, x.Message, x.MessageID)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func register(newfunc commandCtor) {
	commandList = append(commandList, newfunc)
}

func findCommand(command string) (a command) {
	for _, item := range commandList {
		a = item()
		if a.commandIdentifier() == command {
			return a
		}
	};
	return nil
}

func executeCommand(update tgbotapi.Update) bool {
	command := findCommand(update.Message.Command())
	if command != nil {
		command.execute(update)
		return true
	}
	return false
}

func getCommands() []commandDescription {
	result := make([]commandDescription, len(commandList))
	for i, x := range commandList {
		result[i] = commandDescription{x().commandIdentifier(), x().commandDescription()}
	}
	return result
}

func heartbeat() {
	time.Sleep(time.Second * 30)
	for true {
		sendToHeartbeatGroup(emoji.Sprintf("%s :heart:", bot.Self.UserName))
		time.Sleep(time.Hour)
	}
}

func (s *setHeartbeatGroup) commandIdentifier() string {
	return "SetHeartbeatGroup"
}

func (s *setHeartbeatGroup) commandDescription() string {
	return "Set Heartbeat Group"
}

func (s *setHeartbeatGroup) execute(update tgbotapi.Update) {
	database.SetHeartbeatGroup(update.Message.Chat.ID)
	SendMessage(update.Message.Chat.ID, "heartbeat group updated", update.Message.MessageID)
}
