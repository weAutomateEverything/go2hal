package callout

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"golang.org/x/net/context"
	"os"

	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sony/gobreaker"
	"github.com/weAutomateEverything/go2hal/alert"
	"github.com/weAutomateEverything/go2hal/gokit"
	"golang.org/x/time/rate"
	"time"
)

type calloutProxy struct {
	namespace string
	logger    log.Logger
}

// NewCalloutProxy will create a HTTP Rest client to easily invoke the Callout service offered by HAL. The HAL Service
// endpoint needs to be set in a Environment Variable named HAL_ENDPOINT
func NewCalloutProxy() Service {
	if getHalUrl() == "" {
		panic("No Alert Endpoint set. Please set the environment variable ALERT_ENDPOINT with the http address of the alert service")
	}
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	return newProxy("", logger)

}

// NewKubernetesCalloutProxy creates a HTTP Rest client to easily call the Callout service by using the Service discovery
// mechaisms build into kubernetes. The funciton does assume that the deployment was done using the templates provided
// within the kuberets folder of the project, and that a service names hal was created. If you application resides in the
// same namespace as HAL, then the namespace can be left as a empty string, else provide the namespace that contains
// the hal deployment.
func NewKubernetesCalloutProxy(namespace string) Service {
	fieldKeys := []string{"method"}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	service := newProxy(namespace, logger)
	service = NewLoggingService(log.With(logger, "component", "alert_proxy"), service)
	service = NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "proxy",
		Subsystem: "callout_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "proxy",
			Subsystem: "callout_service",
			Name:      "error_count",
			Help:      "Number of errors.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "proxy",
			Subsystem: "callout_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), service)

	return service
}

func newProxy(namespace string, logger log.Logger) Service {

	return &calloutProxy{
		namespace: namespace,
		logger:    logger,
	}
}

func (s *calloutProxy) InvokeCallout(ctx context.Context, chatId uint32, title, message string, variables map[string]string) error {

	callout := makeCalloutHttpProxy(s.namespace, chatId, s.logger)
	callout = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(callout)
	callout = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 10))(callout)

	_, err := callout(ctx, &SendCalloutRequest{Message: message, Title: title})
	return err
}

func (calloutProxy) getFirstCall(ctx context.Context) (name string, number string, err error) {
	panic("Not implemented")
}

func getHalUrl() string {
	return os.Getenv("HAL_ENDPOINT")
}

func makeCalloutHttpProxy(namespace string, chatid uint32, logger log.Logger) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		alert.GetURL(namespace, fmt.Sprintf("callout/%v", chatid)),
		gokit.EncodeJsonRequest,
		gokit.DecodeResponse,
		gokit.GetClientOpts(logger)...,
	).Endpoint()
}
