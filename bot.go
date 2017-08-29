package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/commands"
)

var bot *tgbotapi.BotAPI

type HalBot struct {
	Running bool
	bot *tgbotapi.BotAPI
}

func StartBot() *HalBot{
	bot, err := tgbotapi.NewBotAPI("341557986:AAEX0Ew-0T-co8KNPJz028yqtdd3i9SVD8I")

	if err != nil {
		log.Panic(err)
	}
	hal := &HalBot{true, bot}
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
					commands.ExecuteCommand(update.Message.Command(), update.Message.CommandArguments())
					continue
				}

				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				sendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
			}
		}
	}()

	return hal

}
func sendMessage(chatId int64, message string, messageId int) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ReplyToMessageID = messageId
	bot.Send(msg)
}
