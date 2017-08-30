package telegram

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
)


type HalBot struct {
	Running bool
	bot *tgbotapi.BotAPI
}

var hal *HalBot;
var bot *tgbotapi.BotAPI
var err error

func init(){
	log.Println("Hal Init")
}

func Start() {
	bot, err = tgbotapi.NewBotAPI("303007671:AAHxx3PI9zU5q-4a43WmQXU3K-A6yGY1DAU")

	if err != nil {
		log.Panic(err)
	}
	hal = &HalBot{true, bot}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go func() {
		for true {
			log.Println("Waiting for messages...")
			updates, err := bot.GetUpdatesChan(u)
			if err != nil {
				log.Panic(err)
			}

			for update := range updates {
				if update.Message == nil {
					continue
				}

				if update.Message.IsCommand(){
					ExecuteCommand(update)
					continue
				}

				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				SendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
			}
		}
	}()
}

func GetBot() *HalBot{
	return hal
}

func SendMessage(chatId int64, message string, messageId int) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyToMessageID = messageId
	bot.Send(msg)
}
