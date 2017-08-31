package main

import (
	"time"
	"log"
	"github.com/zamedic/go2hal/telegram"
)

func main() {
	log.Println("Starting HAL")
	log.Println("-------------")
	log.Println("All systems GO!")

	hal := telegram.GetBot()
	for true{
		time.Sleep(time.Second * 5)
		log.Printf("Heartbeat...  [%s]",hal.Running)
	}
}
