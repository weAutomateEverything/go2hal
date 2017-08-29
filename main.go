package main

import (
	"fmt"
	"time"
	"log"
)

func main() {
	fmt.Println("HAL Starting, initialising DB")
	StartDB()
	fmt.Println("Loading rest")
	StartRest()
	fmt.Println("Loading telegram framework")
	StartBot()

	for true{
		time.Sleep(time.Second * 5)
		log.Println("Heartbeat... ")
	}

}
