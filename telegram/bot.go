package telegram

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/zamedic/go2hal/database"
	"time"
)

/**
Structure to describe the state of the bot
 */
type HalBot struct {
	Running bool
	bot *tgbotapi.BotAPI
}

var hal *HalBot;
var bot *tgbotapi.BotAPI
var err error

func init() {
	go func() {
		findFreeBot()
	}()
}

func findFreeBot(){
	for true {
		bots := database.ListBots()
		if bots !=  nil{
			for _, botkey := range bots {
				useBot(botkey)
			}
		}
		log.Println("Looped through all bots. Sleeping for 2 minutes.")
		time.Sleep(time.Minute * 2)
	}
}

func useBot(botkey string){
	bot, err = tgbotapi.NewBotAPI(botkey)

	if err != nil {
		log.Panic(err)
	}
	hal = &HalBot{true, bot}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60


		for true {
			log.Println("Waiting for messages...")
			updates, err := bot.GetUpdatesChan(u)
			if err != nil {
				log.Println("Releasing bot ",bot.Self.UserName)
				log.Panic(err)
				hal.Running = false
				hal.bot = nil
				return
			}

			for update := range updates {
				if update.Message == nil {
					continue
				}

				if update.Message.IsCommand(){
					executeCommand(update)
					continue
				}

				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
				SendMessage(update.Message.Chat.ID, update.Message.Text, update.Message.MessageID)
			}
		}

}

/**
Returns a handler back to the bot
 */
func GetBot() *HalBot{
	return hal
}

/**
Sends a test message to the chat id.
 */
func SendMessage(chatID int64, message string, messageID int) (err error){
	msg := tgbotapi.NewMessage(chatID, message)
	if messageID != 0 {
		msg.ReplyToMessageID = messageID
	}
	result, err := bot.Send(msg)

	if err != nil {
		log.Panic(err)
	}

	log.Println(result)

	return nil
}
