package jira

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

func (s *instrumentingService) CreateJira(ctx context.Context, title, description string, name string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "CreateJira").Add(1)
		s.requestLatency.With("method", "CreateJira").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.CreateJira(ctx, title, description, name)
}
