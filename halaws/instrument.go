package halaws

import (
	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s instrumentingService) SendAlert(ctx context.Context, destination string, name string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendAlert").Add(1)
		s.requestLatency.With("method", "SendAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.SendAlert(ctx, destination, name)
}
