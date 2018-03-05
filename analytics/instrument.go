package analytics

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

func (s instrumentingService) SendAnalyticsAlert(ctx context.Context, message string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_chef_analytics").Add(1)
		s.requestLatency.With("method", "send_chef_analytics").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.SendAnalyticsAlert(ctx, message)
}
