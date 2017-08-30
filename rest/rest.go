package rest

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func Start(){
	r := mux.NewRouter()
	r.HandleFunc("/alert",AlertHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":8000", r))
	}()
}

