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

	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

type alertKubernetesProxy struct {
	namespace string
	logger    log.Logger

	sendErrorEndpoint      endpoint.Endpoint
	sendErrorImageEndpoint endpoint.Endpoint
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

	alertError := makeAlertSendErrorHTTPProxy(namespace, logger)
	alertError = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertError)
	alertError = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertError)

	alertImageError := makeAlertSendErrorImageHTTPProxy(namespace, logger)
	alertImageError = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertError)
	alertImageError = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertError)

	return &alertKubernetesProxy{
		namespace:              namespace,
		logger:                 logger,
		sendErrorEndpoint:      alertError,
		sendErrorImageEndpoint: alertImageError,
	}

}
func (s *alertKubernetesProxy) SendAlert(ctx context.Context, chatId uint32, message string) error {
	alert := makeAlertHTTPProxy(s.namespace, chatId, s.logger)
	alert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alert)
	alert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alert)
	_, err := alert(getContext(ctx), message)
	return err
}

func (s *alertKubernetesProxy) SendImageToAlertGroup(ctx context.Context, chatId uint32, image []byte) error {

	alertImage := makeAlertSendImageToAlertGroupHTTPProxy(s.namespace, chatId, s.logger)
	alertImage = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertImage)
	alertImage = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(alertImage)
	_, err := alertImage(context.WithValue(getContext(ctx), "CHAT-ID", chatId), image)
	return err
}

func (s *alertKubernetesProxy) SendError(ctx context.Context, err error) error {
	_, e := s.sendErrorEndpoint(getContext(ctx), err)
	return e
}

func (s *alertKubernetesProxy) SendErrorImage(ctx context.Context, image []byte) error {
	_, e := s.sendErrorImageEndpoint(getContext(ctx), image)
	return e
}

func makeAlertHTTPProxy(namespace string, chatId uint32, logger log.Logger) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		GetURL(namespace, fmt.Sprintf("alert/%v", chatId)),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()

}

func makeAlertSendImageToAlertGroupHTTPProxy(namespace string, chatId uint32, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, fmt.Sprintf("alert/%v/image", chatId)),
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

func makeAlertSendErrorImageHTTPProxy(namespace string, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		GetURL(namespace, "alert/error/image"),
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
