package halaws

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

func (s instrumentingService) SendAlert(ctx context.Context, destination string, name string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendAlert").Add(1)
		s.requestLatency.With("method", "SendAlert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCounter.With("method", "SendAlert").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlert(ctx, destination, name)
}
