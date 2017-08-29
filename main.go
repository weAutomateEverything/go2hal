package main

import (
	"fmt"
	"time"
	"log"
	"github.com/zamedic/go2hal/database"
)

func main() {
	fmt.Println("HAL Starting, initialising DB")
	database.StartDB()
	fmt.Println("Loading rest")
	StartRest()
	fmt.Println("Loading telegram framework")
	StartBot()

	for true{
		time.Sleep(time.Second * 5)
		log.Println("Heartbeat... ")
	}

}
