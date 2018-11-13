package go2halStartup

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/analytics"
	"github.com/weAutomateEverything/go2hal/appdynamics"
	"github.com/weAutomateEverything/go2hal/auth"
	"github.com/weAutomateEverything/go2hal/callout"
	"github.com/weAutomateEverything/go2hal/chef"
	"github.com/weAutomateEverything/go2hal/database"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/github"
	"github.com/weAutomateEverything/go2hal/halaws"
	"github.com/weAutomateEverything/go2hal/httpSmoke"
	"github.com/weAutomateEverything/go2hal/jira"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"github.com/weAutomateEverything/go2hal/remoteTelegramCommands"
	"github.com/weAutomateEverything/go2hal/seleniumTests"
	"github.com/weAutomateEverything/go2hal/sensu"
	"github.com/weAutomateEverything/go2hal/snmp"
	"github.com/weAutomateEverything/go2hal/ssh"
	"github.com/weAutomateEverything/go2hal/telegram"
	"github.com/weAutomateEverything/go2hal/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/gorilla/handlers"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/weAutomateEverything/go2hal/grafana"
	"github.com/weAutomateEverything/go2hal/prometheus"
)

type Go2Hal struct {
	Logger     log.Logger
	HTTPLogger log.Logger

	Mux *http.ServeMux

	Err chan error

	DB *mgo.Database

	TelegramStore         telegram.Store
	AppdynamicsStore      appdynamics.Store
	ChefStore             chef.Store
	SSHStore              ssh.Store
	UserStore             user.Store
	SeleniumStore         seleniumTests.Store
	HTTPStore             httpSmoke.Store
	MachineLearningStore  machineLearning.Store
	DefaultFirstcallStore firstCall.Store

	AuthService             auth.Service
	MachineLearningService  machineLearning.Service
	TelegramService         telegram.Service
	AlertService            alert.Service
	JiraService             jira.Service
	AnalticsServics         analytics.Service
	SSHService              ssh.Service
	AppdynamicsService      appdynamics.Service
	SNMPService             snmp.Service
	AWSService              halaws.Service
	DefaultFirstcallService firstCall.CalloutFunction
	FirstCallService        firstCall.Service
	CalloutService          callout.Service
	ChefService             chef.Service
	SensuService            sensu.Service
	UserService             user.Service
	GithubService           github.Service
	SeleniumService         seleniumTests.Service
	HTTPService             httpSmoke.Service
	MqService               appdynamics.MqService
	grafanaService          grafana.Service
	prometheusService       prometheus.Service

	RemoteTelegramCommandService remoteTelegramCommands.RemoteCommandServer
	AppDynamics                  bool
}
type GO2HAL interface {
	Start()
}

func (go2hal *Go2Hal) Start() {
	grpc := grpc.NewServer()
	remoteTelegramCommands.RegisterRemoteCommandServer(grpc, go2hal.RemoteTelegramCommandService)
	reflection.Register(grpc)

	errs := make(chan error, 2)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		go2hal.Logger.Log("transport", "grpc", "address", ":8080", "error", err)
		errs <- err
		panic(err)
	}
	go func() {
		go2hal.Logger.Log("transport", "http", "address", ":8000", "msg", "listening")
		errs <- http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, panicHandler{accessControl(go2hal.Mux, go2hal.Logger), go2hal.JiraService, go2hal.AlertService}))
	}()
	go func() {
		go2hal.Logger.Log("transport", "http", "address", ":8080", "msg", "listening")
		errs <- grpc.Serve(ln)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	go2hal.Logger.Log("terminated", <-go2hal.Err)
}

func NewGo2Hal() Go2Hal {

	go2hal := Go2Hal{}

	go2hal.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	go2hal.Logger = level.NewFilter(go2hal.Logger, level.AllowAll())
	go2hal.Logger = log.With(go2hal.Logger, "ts", log.DefaultTimestamp)

	go2hal.DB = database.NewConnection()

	//Stores
	go2hal.TelegramStore = telegram.NewMongoStore(go2hal.DB)
	go2hal.AppdynamicsStore = appdynamics.NewMongoStore(go2hal.DB)
	go2hal.ChefStore = chef.NewMongoStore(go2hal.DB)
	go2hal.SSHStore = ssh.NewMongoStore(go2hal.DB)
	go2hal.UserStore = user.NewMongoStore(go2hal.DB)
	go2hal.SeleniumStore = seleniumTests.NewMongoStore(go2hal.DB)
	go2hal.HTTPStore = httpSmoke.NewMongoStore(go2hal.DB)
	go2hal.MachineLearningStore = machineLearning.NewMongoStore(go2hal.DB)
	go2hal.DefaultFirstcallStore = firstCall.NewMongoStore(go2hal.DB)

	fieldKeys := []string{"method"}

	//Services

	// Used to save all the requests and responses for later machine learning.
	go2hal.MachineLearningService = machineLearning.NewService(go2hal.MachineLearningStore)

	//A auth service for our bank. If you dont want auth use auth.alwaysTrueAuthService()
	go2hal.AuthService = auth.NewAlwaysTrustEveryoneAuthService()

	go2hal.TelegramService = telegram.NewService(go2hal.TelegramStore, go2hal.AuthService)
	go2hal.TelegramService = telegram.NewLoggingService(log.With(go2hal.Logger, "component", "telegram"), go2hal.TelegramService)
	go2hal.TelegramService = telegram.NewMachineLearning(go2hal.MachineLearningService, go2hal.TelegramService)
	go2hal.TelegramService = telegram.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "telgram_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"method", "chat"}),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "telgram_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, []string{"method", "chat"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "telegram_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "chat"}), go2hal.TelegramService)

	go2hal.AlertService = alert.NewService(go2hal.TelegramService, go2hal.TelegramStore)
	go2hal.AlertService = alert.NewLoggingService(log.With(go2hal.Logger, "component", "alert"), go2hal.AlertService)
	go2hal.AlertService = alert.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "alert_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"method", "chat"}),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "alert_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, []string{"method", "chat"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "alert_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method", "chat"}), go2hal.AlertService)

	go2hal.JiraService = jira.NewService(go2hal.AlertService, go2hal.UserStore)
	go2hal.JiraService = jira.NewLoggingService(log.With(go2hal.Logger, "component", "jira"), go2hal.JiraService)
	go2hal.JiraService = jira.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.JiraService)

	go2hal.AnalticsServics = analytics.NewService(go2hal.AlertService, go2hal.ChefStore)
	go2hal.AnalticsServics = analytics.NewLoggingService(log.With(go2hal.Logger, "component", "chef_audir"), go2hal.AnalticsServics)
	go2hal.AnalticsServics = analytics.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.AnalticsServics)

	go2hal.SSHService = ssh.NewService(go2hal.AlertService, go2hal.SSHStore)
	go2hal.SSHService = ssh.NewLoggingService(log.With(go2hal.Logger, "component", "ssh"), go2hal.SSHService)
	go2hal.SSHService = ssh.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.SSHService)

	go2hal.MqService = appdynamics.NewMqSercvice(go2hal.AlertService, go2hal.AppdynamicsStore)
	go2hal.MqService = appdynamics.NewMqLoggingService(log.With(go2hal.Logger, "component", "appdynamics_mq"), go2hal.MqService)
	go2hal.MqService = appdynamics.NewMqInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "appdynamics_mq",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "appdynamics_mq",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), go2hal.MqService)

	go2hal.SNMPService = snmp.NewService(go2hal.AlertService, go2hal.MachineLearningService)
	go2hal.SNMPService = snmp.NewLoggingService(log.With(go2hal.Logger, "component", "snmp"), go2hal.SNMPService)
	go2hal.SNMPService = snmp.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.SNMPService)

	go2hal.AWSService = halaws.NewService(go2hal.AlertService)
	go2hal.AWSService = halaws.NewLoggingService(log.With(go2hal.Logger, "component", "halaws"), go2hal.AWSService)
	go2hal.AWSService = halaws.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.AWSService)

	go2hal.DefaultFirstcallService = firstCall.NewDefaultFirstcallService(go2hal.DefaultFirstcallStore, go2hal.AlertService)
	go2hal.FirstCallService = firstCall.NewCalloutService()
	go2hal.FirstCallService.AddCalloutFunc(go2hal.DefaultFirstcallService)

	go2hal.CalloutService = callout.NewService(go2hal.AlertService, go2hal.FirstCallService, go2hal.SNMPService, go2hal.JiraService, go2hal.AWSService)
	go2hal.CalloutService = callout.NewLoggingService(log.With(go2hal.Logger, "component", "callout"), go2hal.CalloutService)
	go2hal.CalloutService = callout.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.CalloutService)

	go2hal.AppdynamicsService = appdynamics.NewService(go2hal.AlertService, go2hal.SSHService, go2hal.AppdynamicsStore, go2hal.CalloutService)
	go2hal.AppdynamicsService = appdynamics.NewLoggingService(log.With(go2hal.Logger, "component", "appdynamics"), go2hal.AppdynamicsService)
	go2hal.AppdynamicsService = appdynamics.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.AppdynamicsService)

	go2hal.ChefService = chef.NewService(go2hal.AlertService, go2hal.ChefStore)
	go2hal.ChefService = chef.NewLoggingService(log.With(go2hal.Logger, "component", "chef"), go2hal.ChefService)
	go2hal.ChefService = chef.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.ChefService)

	go2hal.SensuService = sensu.NewService(go2hal.AlertService)
	go2hal.SensuService = sensu.NewLoggingService(log.With(go2hal.Logger, "component", "sensu"), go2hal.SensuService)
	go2hal.SensuService = sensu.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.SensuService)

	go2hal.UserService = user.NewService(go2hal.UserStore)
	go2hal.UserService = user.NewLoggingService(log.With(go2hal.Logger, "component", "userService"), go2hal.UserService)
	go2hal.UserService = user.NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
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
		}, fieldKeys), go2hal.UserService)

	go2hal.RemoteTelegramCommandService = remoteTelegramCommands.NewService(go2hal.TelegramService)

	go2hal.GithubService = github.NewService(go2hal.AlertService)

	go2hal.SeleniumService = seleniumTests.NewService(go2hal.SeleniumStore, go2hal.AlertService, go2hal.CalloutService)
	go2hal.HTTPService = httpSmoke.NewService(go2hal.AlertService, go2hal.HTTPStore, go2hal.CalloutService,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "service",
			Subsystem: "http",
			Name:      "test_count",
			Help:      "Number of http tests executed.",
		}, []string{"endpoint"}),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "service",
			Subsystem: "http",
			Name:      "test_error_count",
			Help:      "number of http tests failed.",
		}, []string{"endpoint"}),
	)

	go2hal.grafanaService = grafana.NewService(go2hal.AlertService)
	go2hal.grafanaService = grafana.NewLoggingService(log.With(go2hal.Logger, "component", "grafana_service"), go2hal.grafanaService)
	go2hal.grafanaService = grafana.NewInstrumentService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "grafana_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "grafana_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "grafana_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		go2hal.grafanaService,
	)

	go2hal.prometheusService = prometheus.NewService(go2hal.AlertService)
	go2hal.prometheusService = prometheus.NewLoggingService(log.With(go2hal.Logger, "component", "prometheus_service"), go2hal.prometheusService)
	go2hal.prometheusService = prometheus.NewInstrumentService(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "prometheus_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "prometheus_service",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "prometheus_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		go2hal.prometheusService,
	)

	//Telegram Commands
	go2hal.TelegramService.RegisterCommand(telegram.NewHelpCommand(go2hal.TelegramService, go2hal.TelegramStore))
	go2hal.TelegramService.RegisterCommand(firstCall.NewWhosOnFirstCallCommand(go2hal.AlertService, go2hal.TelegramService,
		go2hal.FirstCallService, go2hal.TelegramStore))
	go2hal.TelegramService.RegisterCommand(httpSmoke.NewQuietHttpAlertCommand(go2hal.TelegramService, go2hal.HTTPService))
	go2hal.TelegramService.RegisterCommand(telegram.NewIDCommand(go2hal.TelegramService, go2hal.TelegramStore))
	go2hal.TelegramService.RegisterCommand(telegram.NewTokenCommand(go2hal.TelegramService, go2hal.TelegramStore))

	go2hal.TelegramService.RegisterCommandLet(telegram.NewTelegramAuthApprovalCommand(go2hal.TelegramService, go2hal.TelegramStore))

	go2hal.HTTPLogger = log.With(go2hal.Logger, "component", "http")

	go2hal.Mux = http.NewServeMux()

	go2hal.Mux.Handle("/api/alert/", alert.MakeHandler(go2hal.AlertService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/chefAudit", analytics.MakeHandler(go2hal.AnalticsServics, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/appdynamics/{chatid:[0-9]+}/queue", appdynamics.MakeMqHandler(go2hal.MqService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/appdynamics/", appdynamics.MakeHandler(go2hal.AppdynamicsService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/chef/", chef.MakeHandler(go2hal.ChefService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/sensu/", sensu.MakeHandler(go2hal.SensuService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/users/", user.MakeHandler(go2hal.UserService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/aws/sendTestAlert", halaws.MakeHandler(go2hal.AWSService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/callout/", callout.MakeHandler(go2hal.CalloutService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/github/", github.MakeHandler(go2hal.GithubService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/telegram/", telegram.MakeHandler(go2hal.TelegramService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/httpEndpoints", httpSmoke.MakeHandler(go2hal.HTTPService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/firstcall/defaultCallout", firstCall.MakeHandler(go2hal.DefaultFirstcallService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/grafana/", grafana.MakeHandler(go2hal.grafanaService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/prometheus/", prometheus.MakeHandler(go2hal.prometheusService, go2hal.HTTPLogger, go2hal.MachineLearningService))
	go2hal.Mux.Handle("/api/metrics", promhttp.Handler())
	go2hal.Mux.Handle("/api/swagger.json", swagger{})

	return go2hal
}

func accessControl(h http.Handler, httpLogger log.Logger) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")

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
				h.alert.SendError(r.Context(), fmt.Errorf("panic detected: %v \n %v", err.Error(), string(debug.Stack())))
			}
		}
	}()
	h.Handler.ServeHTTP(w, r)
}

type swagger struct {
	http.Handler
}

func (h swagger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile("swagger.json")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(200)
		w.Write(b)
	}
}
