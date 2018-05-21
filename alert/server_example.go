package alert

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/weAutomateEverything/go2hal/auth"
	"github.com/weAutomateEverything/go2hal/database"
	"github.com/weAutomateEverything/go2hal/telegram"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func server() {
	db := database.NewConnection()

	//Stores
	telegramStore := telegram.NewMongoStore(db)

	authService := auth.NewAlwaysTrustEveryoneAuthService()

	telegramService := telegram.NewService(telegramStore, authService)
	alertService := NewService(telegramService, telegramStore)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)
	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert/", MakeHandler(alertService, httpLogger, nil))
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
