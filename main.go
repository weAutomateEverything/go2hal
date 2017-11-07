package main

import (
	"time"
	"log"
	"github.com/zamedic/go2hal/rest"
	"github.com/zamedic/go2hal/service"
)

func main() {
	log.Println("Starting HAL")
	log.Println("-------------")
	log.Println("All systems GO!")

	hal := service.GetBot()
	router := rest.Router()

	for true{
		time.Sleep(time.Minute * 5)
		log.Printf("Heartbeat...  Bot: [%v], router: [%v]",hal.Running,router.Mux)
	}
}
