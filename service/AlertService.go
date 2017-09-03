package service

import (
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/telegram"
)

func SendAlert(message string) error {
	alertGroup, err := database.AlertGroup()
	if err != nil{
		return err
	}
	err = telegram.SendMessage(alertGroup, message, 0)
	return err
}
