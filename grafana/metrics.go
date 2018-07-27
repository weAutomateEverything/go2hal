package grafana

import (
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	errorCounter   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, errorCounter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		errorCounter:   errorCounter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) sendGrafanaAlert(ctx context.Context, chat uint32, body string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "sendGrafanaAlert").Add(1)
		if err != nil {
			s.errorCounter.With("method", "sendGrafanaAlert").Add(1)
		}
		s.requestLatency.With("method", "sendGrafanaAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.sendGrafanaAlert(ctx, chat, body)
}
