package service

import (
	"github.com/zamedic/go2hal/database"
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"runtime"
	"io/ioutil"
	"os"
	"math/rand"
)

func init(){
	log.Println("Initializing Set Group Command")
	register(func() command {
		return &setGroup{}
	})
	register(func() command {
		return &setNonTechnicalGroup{}
	})
}

//SendAlert will send the alert message as defined in the alert group
func SendAlert(message string) error {
	alertGroup, err := database.AlertGroup()
	if err != nil{
		return err
	}
	err = SendMessage(alertGroup, message, 0)
	return err
}

//SendAlert will send the alert message as defined in the alert group
func SendNonTechnicalAlert(message string) error {
	alertGroup, err := database.NonTechnicalGroup()
	if err != nil{
		return err
	}
	err = SendMessage(alertGroup, message, 0)
	return err
}

/*
SendImageToAlertGroup will send the image to the alert group
 */
func SendImageToAlertGroup(image []byte) error {


	alertGroup, err := database.AlertGroup()
	if err != nil {
		SendError(err)
		return err
	}

	return sendImageToGroup(image, alertGroup)
}

/*
SendImageToHeartbeatGroup will send the image to the heartbeat group
 */
func SendImageToHeartbeatGroup(image []byte) error {
	group, err := database.HeartbeatGroup()
	if err != nil {
		SendError(err)
		return err
	}

	return sendImageToGroup(image, group)

}

func sendImageToGroup(image []byte, group int64) error{
	path := ""

	if runtime.GOOS == "windows" {
		path = "c:/temp/"
	} else {
		path = "/tmp/"
	}
	path = path + string(rand.Int()) + ".png"
	err := ioutil.WriteFile(path, image, os.ModePerm)

	if err != nil {
		SendError(err)
		return err
	}


	msg := tgbotapi.NewPhotoUpload(group, path)
	if bot != nil{
		_, err = bot.Send(msg)
		if err != nil {
			SendError(err)
			return err
		}
	}

	return nil
}

type setGroup struct {
}

func (s *setGroup) commandIdentifier() string {
	return "SetGroup"
}

func (s *setGroup) commandDescription() string {
	return "Set Alert Group"
}

func (s *setGroup) execute(update tgbotapi.Update){
	database.SetAlertGroup(update.Message.Chat.ID)
	SendMessage(update.Message.Chat.ID,"group updated", update.Message.MessageID)
}

type setNonTechnicalGroup struct {

}

func (s *setNonTechnicalGroup) commandIdentifier() string {
	return "SetNonTechGroup"
}

func (s *setNonTechnicalGroup) commandDescription() string {
	return "Set Non Technical Alert Group"
}

func (s *setNonTechnicalGroup) execute(update tgbotapi.Update){
	database.SetNonTechnicalGroup(update.Message.Chat.ID)
	SendMessage(update.Message.Chat.ID,"non technical group updated", update.Message.MessageID)
}