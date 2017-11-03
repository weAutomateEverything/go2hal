package telegram

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/database"
	"time"
)


//HalBot Structure to describe the state of the bot
type HalBot struct {
	Running bool
	bot     *tgbotapi.BotAPI
}

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
}

/**
Returns a handler back to the bot
 */
func GetBot() *HalBot {
	return hal
}

/**
Sends a test message to the chat id.
 */
func SendMessage(chatID int64, message string, messageID int) (err error) {
	if !hal.Running {
		log.Println("Unable to send message as no bot is connected. ")
		return nil
	}
	log.Printf("Sending Message %s",message)
	if (!hal.Running) {
		database.AddMessageToQueue(message, chatID, messageID)
	}
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
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
		log.Println("Looped through all bots. Sleeping for 2 minutes.")
		time.Sleep(time.Minute * 2)
	}
}

func useBot(botkey string) {
	bot, err = tgbotapi.NewBotAPI(botkey)

	if err != nil {
		hal.Running = false
		log.Printf("Error getting bot token: %s",err.Error())
		return
	}
	hal.Running = true
	hal.bot = bot

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for true {
		log.Println("Waiting for messages...")
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
