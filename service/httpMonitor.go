package service

import (
	"github.com/zamedic/go2hal/database"
	"time"
	"net/http"
	"fmt"
	"log"
)

//HttpMonitor is the current status of the monitor
type HttpMonitor struct {
	running bool
}

var h *HttpMonitor

func init(){
	h = &HttpMonitor{}
	go func() {
		monitorEndpoints()
	}()
}

func monitorEndpoints(){
	log.Println("Starting HTTP Endpoint monitor")
	h.running = true
	for true {
		endpoints := database.GetHtmlEndpoints()
		if endpoints != nil {
			for _, endpoint := range endpoints {
				response, err := http.Get(endpoint.Endpoint)
				if err != nil {
					SendAlert(fmt.Sprintf("*HTTP Alert*\nName: %s \nEndpoint: %s \nError: %s",endpoint.Name,
						endpoint.Endpoint,err.Error()))
					continue
				}
				if response.StatusCode != 200 {
					SendAlert(fmt.Sprintf("*HTTP Alert*\nName: %s \nEndpoint: %s \nDid not receive a 200 success " +
						"response code. Recieved %d response code.",endpoint.Name,endpoint.Endpoint,
						response.StatusCode))
					continue
				}
			}
		}
		time.Sleep(time.Minute * 2)
	}
}
