package main

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func StartRest(){
	r := mux.NewRouter()
	r.HandleFunc("/alert",AlertHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8000", r))
	}()
}

