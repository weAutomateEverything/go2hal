package main

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/analytics"
	"github.com/weAutomateEverything/go2hal/appdynamics"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/chef"
	"github.com/weAutomateEverything/go2hal/database"
	"github.com/weAutomateEverything/go2hal/jira"
	"github.com/weAutomateEverything/go2hal/remoteTelegramCommands"
	"github.com/weAutomateEverything/go2hal/seleniumTests"
	"github.com/weAutomateEverything/go2hal/sensu"
	"github.com/weAutomateEverything/go2hal/skynet"
	"github.com/weAutomateEverything/go2hal/snmp"
	ssh2 "github.com/weAutomateEverything/go2hal/ssh"
	"github.com/weAutomateEverything/go2hal/telegram"
	"github.com/weAutomateEverything/go2hal/user"
	"google.golang.org/grpc/reflection"

	"github.com/weAutomateEverything/bankCallout"
	"github.com/weAutomateEverything/bankldapService"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/github"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/httpSmoke"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	db := database.NewConnection()

	//Stores
	alertStore := alert.NewStore(db)
	telegramStore := telegram.NewMongoStore(db)
	appdynamicsStore := appdynamics.NewMongoStore(db)
	chefStore := chef.NewMongoStore(db)
	sshStore := ssh2.NewMongoStore(db)
	userStore := user.NewMongoStore(db)
	seleniumStore := seleniumTests.NewMongoStore(db)
	httpStore := httpSmoke.NewMongoStore(db)
	machingLearningStore := machineLearning.NewMongoStore(db)

	//A datastore to store the auth info for our bank - not required if you dont want auth
	bankLdapStore := bankldapService.NewMongoStore(db)

	fieldKeys := []string{"method"}

	//Services

	// Used to save all the requests and responses for later machine learning.
	machineLearningService := machineLearning.NewService(machingLearningStore)

	//A auth service for our bank. If you dont want auth use auth.alwaysTrueAuthService()
	authService := bankldapService.NewService(bankLdapStore)

	telegramService := telegram.NewService(telegramStore, authService)
	telegramService = telegram.NewLoggingService(log.With(logger, "component", "telegram"), telegramService)
	telegramService = telegram.NewMachineLearning(machineLearningService, telegramService)
	telegramService = telegram.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "telgram_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "telgram_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "telegram_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), telegramService)

	alertService := alert.NewService(telegramService, alertStore)
	alertService = alert.NewLoggingService(log.With(logger, "component", "alert"), alertService)
	alertService = alert.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "alert_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "alert_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "alert_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), alertService)

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

	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if ok {
				jiraService.CreateJira(context.TODO(), fmt.Sprintf("HAL Panic - %v", err.Error()), string(debug.Stack()), getTechnicalUser())
			}
			panic(r)
		}
	}()

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

	snmpService := snmp.NewService(alertService, machineLearningService)
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

	aws := halaws.NewService(alertService)
	aws = halaws.NewLoggingService(log.With(logger, "component", "halaws"), aws)
	aws = halaws.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "halaws",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "halaws",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "halaws",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), aws)

	firstcallService := bankCallout.NewService()

	calloutService := callout.NewService(alertService, firstcallService, snmpService, jiraService, aws)
	calloutService = callout.NewLoggingService(log.With(logger, "component", "callout"), calloutService)
	calloutService = callout.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "callout",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "callout",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
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

	userService := user.NewService(userStore)
	userService = user.NewLoggingService(log.With(logger, "component", "userService"), userService)
	userService = user.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), userService)

	remoteTelegramCommand := remoteTelegramCommands.NewService(telegramService)

	githubService := github.NewService(alertService)

	_ = seleniumTests.NewService(seleniumStore, alertService, calloutService)
	httpService := httpSmoke.NewService(alertService, httpStore, calloutService)

	//Telegram Commands
	telegramService.RegisterCommand(alert.NewSetGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetNonTechnicalGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(alert.NewSetHeartbeatGroupCommand(telegramService, alertStore))
	telegramService.RegisterCommand(telegram.NewHelpCommand(telegramService))
	telegramService.RegisterCommand(firstCall.NewWhosOnFirstCallCommand(alertService, telegramService, firstcallService))
	telegramService.RegisterCommand(skynet.NewRebuildCHefNodeCommand(telegramStore, chefStore, telegramService,
		alertService))
	telegramService.RegisterCommand(skynet.NewRebuildNodeCommand(alertService, skynetService))
	telegramService.RegisterCommand(httpSmoke.NewQuietHttpAlertCommand(telegramService, httpService))
	telegramService.RegisterCommand(bankldapService.NewRegisterCommand(telegramService, bankLdapStore))
	telegramService.RegisterCommand(bankldapService.NewTokenCommand(telegramService, bankLdapStore))

	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeEnvironmentReplyCommandlet(telegramService,
		skynetService, chefService))
	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeExecute(skynetService, alertService))
	telegramService.RegisterCommandLet(skynet.NewRebuildChefNodeRecipeReplyCommandlet(chefStore, alertService,
		telegramService))

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/alert/", alert.MakeHandler(alertService, httpLogger, machineLearningService))
	mux.Handle("/chefAudit", analytics.MakeHandler(analyticsService, httpLogger, machineLearningService))
	mux.Handle("/appdynamics/", appdynamics.MakeHandler(appdynamicsService, httpLogger, machineLearningService))
	mux.Handle("/delivery", chef.MakeHandler(chefService, httpLogger, machineLearningService))
	mux.Handle("/skynet/", skynet.MakeHandler(skynetService, httpLogger, machineLearningService))
	mux.Handle("/sensu", sensu.MakeHandler(sensuService, httpLogger, machineLearningService))
	mux.Handle("/users/", user.MakeHandler(userService, httpLogger, machineLearningService))
	mux.Handle("/aws/sendTestAlert", halaws.MakeHandler(aws, httpLogger, machineLearningService))
	mux.Handle("/callout/", callout.MakeHandler(calloutService, httpLogger, machineLearningService))
	mux.Handle("/github/", github.MakeHandler(githubService, httpLogger, machineLearningService))

	http.Handle("/", panicHandler{accessControl(mux), jiraService, alertService})
	http.Handle("/metrics", promhttp.Handler())

	grpc := grpc.NewServer()
	remoteTelegramCommands.RegisterRemoteCommandServer(grpc, remoteTelegramCommand)
	reflection.Register(grpc)

	errs := make(chan error, 2)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Log("transport", "grpc", "address", ":8080", "error", err)
		errs <- err
		return
	}
	go func() {
		logger.Log("transport", "http", "address", ":8000", "msg", "listening")
		errs <- http.ListenAndServe(":8000", nil)
	}()
	go func() {
		logger.Log("transport", "http", "address", ":8080", "msg", "listening")
		errs <- grpc.Serve(ln)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)

}
func getTechnicalUser() string {
	return os.Getenv("TECH_USER")
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

type panicHandler struct {
	http.Handler
	jira  jira.Service
	alert alert.Service
}

func (h panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {

			err, ok := err.(error)
			if ok {
				h.alert.SendError(context.TODO(), fmt.Errorf("panic detected: %v \n %v", err.Error(), string(debug.Stack())))
				h.jira.CreateJira(context.TODO(), fmt.Sprintf("HAL Panic - %v", err.Error()), string(debug.Stack()), getTechnicalUser())
			}
		}
	}()
	h.Handler.ServeHTTP(w, r)
}
