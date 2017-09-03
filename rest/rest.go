package rest

import(
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

type RouterObject struct {
	Mux *mux.Router
}

var router *RouterObject

func init(){
	router = &RouterObject{}
	go func() {
		log.Println("Starting HTTP Server...")
		log.Fatal(http.ListenAndServe(":8000", getRouter()))
	}()
}

func getRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/alert",alertHandler)
	r.HandleFunc("/alert/aws/container",handleEc2ContainerAlert)
	r.HandleFunc("/github",handleGithubMessage)
	r.HandleFunc("/status",status)
	r.HandleFunc("/bot",addBot).Methods("POST")
	r.HandleFunc("/bot",botStatus).Methods("GET")
	router.Mux = r
	return r
}

func Router() (*RouterObject) {
	return router
}

