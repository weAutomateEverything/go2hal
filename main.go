package main

import (
	"time"
	"log"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/telegram"
	"github.com/zamedic/go2hal/rest"
)

func main() {
	log.Println("Starting HAL")
	log.Println("-------------")
	log.Println("Starting Database Connection")
	database.Start()
	log.Println("Starting Telegram Connection")
	telegram.Start()
	log.Println("Starting Resful Service")
	rest.Start()
	log.Println("-------------")
	log.Println("All systems GO!")

	hal := telegram.GetBot()
	for true{
		time.Sleep(time.Second * 5)
		log.Printf("Heartbeat...  [%s]",hal.Running)
	}
}
