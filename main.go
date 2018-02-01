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
	ssh2 "github.com/zamedic/go2hal/ssh"
	"github.com/zamedic/go2hal/jira"
	"github.com/zamedic/go2hal/user"
	"github.com/zamedic/go2hal/skynet"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os/signal"
	"syscall"
	"fmt"
	"github.com/zamedic/go2hal/sensu"
	"github.com/go-kit/kit/log/level"
	"github.com/zamedic/go2hal/seleniumTests"
	http2 "github.com/zamedic/go2hal/http"
)

func main() {

	db := database.NewConnection()

	//Stores
	alertStore := alert.NewStore(db)
	telegramStore := telegram.NewMongoStore(db)
	appdynamicsStore := appdynamics.NewMongoStore(db)
	chefStore := chef.NewMongoStore(db)
	sshStore := ssh2.NewMongoStore(db)
	userStore := user.NewMongoStore(db)
	seleniumStore := seleniumTests.NewMongoStore(db)
	httpStore := http2.NewMongoStore(db)

	fieldKeys := []string{"method"}

	//Services
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	telegramService := telegram.NewService(telegramStore)

	alertService := alert.NewService(telegramService, alertStore)
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

	sshService := ssh2.NewService(alertService, sshStore)
	sshService = ssh2.NewLoggingService(log.With(logger, "component", "ssh"), sshService)
	sshService = ssh2.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "ssh",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "ssh",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), sshService)

	appdynamicsService := appdynamics.NewService(alertService, sshService, appdynamicsStore)
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
	snmpService = snmp.NewLoggingService(log.With(logger, "component", "snmp"), snmpService)
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

	jiraService := jira.NewService(alertService, userStore)
	jiraService = jira.NewLoggingService(log.With(logger, "component", "jira"), jiraService)
	jiraService = jira.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "jira",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "jira",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), jiraService)

	calloutService := callout.NewService(alertService, snmpService, jiraService)
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

	chefService := chef.NewService(alertService, chefStore)
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

	skynetService := skynet.NewService(alertService, chefStore, calloutService)
	skynetService = skynet.NewLoggingService(log.With(logger, "component", "skynet"), skynetService)
	skynetService = skynet.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "skynet",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "skynet",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), skynetService)

	sensuService := sensu.NewService(alertService)
	sensuService = sensu.NewLoggingService(log.With(logger, "component", "sensu"), sensuService)
	sensuService = sensu.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "sensu",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "sensu",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), sensuService)

	_ = seleniumTests.NewService(seleniumStore,alertService,calloutService)
	_ = http2.NewService(alertService,httpStore,calloutService)

	//Telegram Commands
	telegramService.RegisterCommand(alert.NewSetGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetNonTechnicalGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetHeartbeatGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(telegram.NewHelpCommand(telegramService))
	telegramService.RegisterCommand(callout.NewWhosOnFirstCallCommand(alertService, telegramService, calloutService))
	telegramService.RegisterCommand(skynet.NewRebuildCHefNodeCommand(telegramStore, chefStore, telegramService,
		alertService))
	telegramService.RegisterCommand(skynet.NewRebuildNodeCommand(alertService, skynetService))

	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeEnvironmentReplyCommandlet(telegramService,
		skynetService, chefService))
	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeExecute(skynetService, alertService))
	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeRecipeReplyCommandlet(chefStore, alertService,
		telegramService))

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert", alert.MakeHandler(alertService, httpLogger))
	mux.Handle("/chefAudit", analytics.MakeHandler(analyticsService, httpLogger))
	mux.Handle("/appdynamics", appdynamics.MakeHandler(appdynamicsService, httpLogger))
	mux.Handle("/delivery", chef.MakeHandler(chefService, httpLogger))
	mux.Handle("/skynet", skynet.MakeHandler(skynetService, httpLogger))
	mux.Handle("/sensu", sensu.MakeHandler(sensuService, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

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
