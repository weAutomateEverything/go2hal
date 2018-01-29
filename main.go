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
	"github.com/zamedic/go2hal/appdynamics"
	"github.com/zamedic/go2hal/callout"
	"github.com/zamedic/go2hal/snmp"
)

func main() {
	db := database.NewConnection()

	//Stores
	alertStore := alert.NewStore(db)
	telegramStore := telegram.NewMongoStore(db)
	appdynamicsStore := appdynamics.NewMongoStore(db)
	chefStore := chef.NewMongoStore(db)

	fieldKeys := []string{"method"}

	//Services
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	telegramService := telegram.NewService(telegramStore)

	alertService := alert.NewService(telegramService)
	alertService = alert.NewLoggingService(log.With(logger, "component", "alert"), alertService)
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
	analyticsService = analytics.NewLoggingService(log.With(logger, "component", "chef_audir"), analyticsService)
	analyticsService = analytics.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "chef_audit",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "chef_audit",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), analyticsService)

	appdynamicsService := appdynamics.NewService(alertService,sshService,appdynamicsStore)
	appdynamicsService = appdynamics.NewLoggingService(log.With(logger, "component", "appdynamics"), appdynamicsService)
	appdynamicsService = appdynamics.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "appdynamics",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "appdynamics",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), appdynamicsService)

	snmpService := snmp.NewService(alertService)
	snmpService = snmp.NewLoggingService(log.With(logger, "component", "snmp"),snmpService)
	snmpService = snmp.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "snmp",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "snmp",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), snmpService)



	calloutService := callout.NewService(alertService,snmpService,jiraService)
	calloutService = callout.NewLoggingService(log.With(logger, "component", "callout"), calloutService)
	calloutService = callout.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "callout",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "callout",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), calloutService)

	chefService := chef.NewService(alertService,chefStore)
	chefService = chef.NewLoggingService(log.With(logger, "component", "chef"), chefService)
	chefService = chef.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "chef",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "chef",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), chefService)


	//Telegram Commands
	telegramService.RegisterCommand(alert.NewSetGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetNonTechnicalGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetHeartbeatGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(telegram.NewHelpCommand(telegramService))
	telegramService.RegisterCommand(callout.NewWhosOnFirstCallCommand(alertService,telegramService,calloutService))

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert", alert.MakeHandler(alertService, httpLogger))
	mux.Handle("/audit", analytics.MakeHandler(analyticsService, httpLogger))
	mux.Handle("/appdynamics",appdynamics.MakeHandler(appdynamicsService,httpLogger))
	mux.Handle("/chef",chef.MakeHandler(chefService,httpLogger))

}
