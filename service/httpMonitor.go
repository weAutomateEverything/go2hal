package service

import (
	"github.com/zamedic/go2hal/database"
	"time"
	"net/http"
	"fmt"
	"log"
	"crypto/tls"
)

//HTTPMonitor is the current status of the monitor
type HTTPMonitor struct {
	running bool
}

var h *HTTPMonitor

func init(){
	h = &HTTPMonitor{}
	go func() {
		monitorEndpoints()
	}()
}

func monitorEndpoints(){
	log.Println("Starting HTTP Endpoint monitor")
	h.running = true
	for true {
		endpoints := database.GetHTMLEndpoints()
		if endpoints != nil {
			for _, endpoint := range endpoints {
				checkHTTP(endpoint)
			}
		}
		time.Sleep(time.Minute * 2)
	}
}

func checkHTTP(endpoint database.HTMLEndpoint){
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(endpoint.Endpoint)
	if err != nil {
		SendAlert(fmt.Sprintf("*HTTP Alert*\nName: %s \nEndpoint: %s \nError: %s",endpoint.Name,
			endpoint.Endpoint,err.Error()))
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		SendAlert(fmt.Sprintf("*HTTP Alert*\nName: %s \nEndpoint: %s \nDid not receive a 200 success " +
			"response code. Recieved %d response code.",endpoint.Name,endpoint.Endpoint,
			response.StatusCode))
	}
}
