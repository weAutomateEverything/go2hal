package main

import (
	"github.com/zamedic/go2hal/telegram"
	"github.com/zamedic/go2hal/alert"
	"github.com/zamedic/go2hal/database"
	"github.com/zamedic/go2hal/analytics"
	"github.com/zamedic/go2hal/chef"
	"net/http"
	"github.com/go-kit/kit/log"
	"os"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	db := database.NewConnection()

	//Stores
	alertStore := alert.NewStore(db)
	telegramStore := telegram.NewMongoStore(db)

	chefStore := chef.NewMongoStore(db)

	fieldKeys := []string{"method"}

	//Services
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)



	telegramService := telegram.NewService(telegramStore)

	alertService := alert.NewService(telegramService)
	alertService = alert.NewLoggingService(log.With(logger,"component","alert"),alertService)
	alertService = alert.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "alert_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "alert_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), alertService)


	analyticsService := analytics.NewService(alertService, chefStore)

	//Telegram Commands
	telegramService.RegisterCommand(alert.NewSetGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetNonTechnicalGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetHeartbeatGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(telegram.NewHelpCommand(telegramService))



	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert", alert.MakeHandler(alertService, httpLogger))

}
