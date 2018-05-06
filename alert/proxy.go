package alert

import (
	"context"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"github.com/weAutomateEverything/go2hal/gokit"
	"golang.org/x/time/rate"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type alertKubernetesProxy struct {
	sendAlertEndpoint                 endpoint.Endpoint
	sendNonTechnicalAlertEndpoint     endpoint.Endpoint
	sendHeartbeatGroupAlertEndpoint   endpoint.Endpoint
	sendImageToAlertGroupEndpoint     endpoint.Endpoint
	sendImageToHeartbeatGroupEndpoint endpoint.Endpoint
	sendErrorEndpoint                 endpoint.Endpoint
	sendKeyboardRecipeAlertEndpoint   endpoint.Endpoint
	sendEnvironmentAlertEndpoint      endpoint.Endpoint
	sendNodesAlertEndpoint            endpoint.Endpoint
}

/*
NewAlertProxy will return an alert service that is actually a HTTP Proxy into the alert service as defined by the ALERT_ENDPOINT
environment variable.

If the environment variable ALERT_ENDPOINT is blank, then a panic will be raised.
*/
func NewAlertProxy() Service {
	if getHalUrl() == "" {
		panic("No Alert Endpoint set. Please set the environment variable ALERT_ENDPOINT with the http address of the alert service")
	}
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	return newProxy("", logger)
}

/*
NewKubernetesAlertProxy will return an alert service that is actually a HTTP Proxy into the kubertes service
*/
func NewKubernetesAlertProxy(namespace string) Service {

	fieldKeys := []string{"method"}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	service := newProxy(namespace, logger)
	service = NewLoggingService(log.With(logger, "component", "alert_proxy"), service)
	service = NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "proxy",
		Subsystem: "alert_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "proxy",
			Subsystem: "alert_service",
			Name:      "error_count",
			Help:      "Number of errors.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "proxy",
			Subsystem: "alert_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), service)

	return service
}
func newProxy(namespace string, logger log.Logger) Service {
	alert := makeAlertHTTPProxy(namespace, logger)
	alert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alert)
	alert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alert)

	alertKeyboardRecipe := makeKeyboardRecipeAlertHTTPProxy(namespace, logger)
	alertKeyboardRecipe = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertKeyboardRecipe)
	alertKeyboardRecipe = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertKeyboardRecipe)

	alertEnvironment := makeEnvironmentAlertHTTPProxy(namespace, logger)
	alertEnvironment = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertEnvironment)
	alertEnvironment = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertEnvironment)

	alertNodes := makeNodesAlertHTTPProxy(namespace, logger)
	alertNodes = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertNodes)
	alertNodes = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertNodes)

	alertImage := makeAlertSendImageToAlertGroupHTTPProxy(namespace, logger)
	alertImage = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertImage)
	alertImage = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertImage)

	nonTechAlert := makeAlertSendNonTechnicalAlertHTTPProxy(namespace, logger)
	nonTechAlert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(nonTechAlert)
	nonTechAlert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(nonTechAlert)

	heartbeatAlert := makeAlertSendHeartbeatGroupAlertHTTPProxy(namespace, logger)
	heartbeatAlert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(heartbeatAlert)
	heartbeatAlert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(heartbeatAlert)

	heartbeatImage := makeAlertSendImageToHeartbeatGroupHTTPProxy(namespace, logger)
	heartbeatImage = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(heartbeatImage)
	heartbeatImage = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(heartbeatImage)

	alertError := makeAlertSendErrorHTTPProxy(namespace, logger)
	alertError = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertError)
	alertError = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertError)

	return &alertKubernetesProxy{sendAlertEndpoint: alert, sendErrorEndpoint: alertError, sendKeyboardRecipeAlertEndpoint: alertKeyboardRecipe, sendEnvironmentAlertEndpoint: alertEnvironment,
		sendNodesAlertEndpoint: alertNodes, sendHeartbeatGroupAlertEndpoint: heartbeatAlert, sendImageToAlertGroupEndpoint: alertImage,
		sendImageToHeartbeatGroupEndpoint: heartbeatImage, sendNonTechnicalAlertEndpoint: nonTechAlert}

}
func (s *alertKubernetesProxy) SendAlert(ctx context.Context, message string) error {
	_, err := s.sendAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendAlertKeyboardRecipe(ctx context.Context, message []string) error {
	_, err := s.sendKeyboardRecipeAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendAlertEnvironment(ctx context.Context, message []string) error {
	_, err := s.sendEnvironmentAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendAlertNodes(ctx context.Context, message []string) error {
	_, err := s.sendNodesAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendNonTechnicalAlert(ctx context.Context, message string) error {
	_, err := s.sendNonTechnicalAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendHeartbeatGroupAlert(ctx context.Context, message string) error {
	_, err := s.sendHeartbeatGroupAlertEndpoint(getContext(ctx), message)
	return err
}
func (s *alertKubernetesProxy) SendImageToAlertGroup(ctx context.Context, image []byte) error {
	_, err := s.sendImageToAlertGroupEndpoint(getContext(ctx), image)
	return err
}
func (s *alertKubernetesProxy) SendImageToHeartbeatGroup(ctx context.Context, image []byte) error {
	_, err := s.sendImageToHeartbeatGroupEndpoint(getContext(ctx), image)
	return err
}
func (s *alertKubernetesProxy) SendError(ctx context.Context, err error) error {
	_, e := s.sendErrorEndpoint(getContext(ctx), err)
	return e
}

func makeAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()

}

func makeAlertSendNonTechnicalAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/business"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}
func makeKeyboardRecipeAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/keyboard/recipe"),
		gokit.EncodeJsonRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()

}
func makeEnvironmentAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/keyboard/environment"),
		gokit.EncodeJsonRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()

}
func makeNodesAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/keyboard/nodes"),
		gokit.EncodeJsonRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()

}
func makeAlertSendHeartbeatGroupAlertHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/heartbeat"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}

func makeAlertSendImageToAlertGroupHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/image"),
		gokit.EncodeToBase64,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}

func makeAlertSendImageToHeartbeatGroupHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/heartbeat/image"),
		gokit.EncodeToBase64,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}

func makeAlertSendErrorHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/error"),
		gokit.EncodeErrorRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}

func GetURL(namespace, uri string) *url.URL {
	u := getHalUrl()
	if u != "" {
		u = u + uri
		ur, err := url.Parse(u)
		if err != nil {
			panic(err)
		}
		return ur
	}
	u = "http://hal"
	if namespace != "" {
		u = u + "." + namespace
	}
	u = u + "/" + uri

	ur, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return ur
}

func getHalUrl() string {
	return os.Getenv("HAL_ENDPOINT")
}

func getContext(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	return ctx
}
