package rest

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func init(){
	go func() {
		log.Fatal(http.ListenAndServe(":8000", getRouter()))
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/alert",alertHandler)
	r.HandleFunc("/status",status)
	return r
}

