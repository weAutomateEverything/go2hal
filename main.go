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

	service.GetBot()
	rest.Router()

	for true{
		time.Sleep(time.Minute * 5)
	}
}
