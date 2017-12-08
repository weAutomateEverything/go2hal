package main

import (
	"time"
	"log"
	"github.com/zamedic/go2hal/rest"
	"github.com/zamedic/go2hal/service"
	"fmt"
	"errors"
	"runtime/debug"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Print(err)
			service.SendError(errors.New(fmt.Sprint(err)))
			service.SendError(errors.New(string(debug.Stack())))
		}
	}()

	log.Println("Starting HAL")
	log.Println("-------------")
	log.Println("All systems GO!")

	service.GetBot()
	rest.Router()

	for true{
		time.Sleep(time.Minute * 5)
	}
}
