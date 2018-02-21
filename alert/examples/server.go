package examples

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/telegram"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func server() {
	db := database.NewConnection()

	//Stores
	alertStore := alert.NewStore(db)
	telegramStore := telegram.NewMongoStore(db)

	telegramService := telegram.NewService(telegramStore)
	alertService := alert.NewService(telegramService, alertStore)

	telegramService.RegisterCommand(alert.NewSetGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetNonTechnicalGroupCommand(telegramService, alertStore))

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert/", alert.MakeHandler(alertService, httpLogger))
	http.Handle("/", accessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", ":8000", "msg", "listening")
		errs <- http.ListenAndServe(":8000", nil)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
