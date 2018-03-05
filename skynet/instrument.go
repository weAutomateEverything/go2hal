package skynet

import (
	"context"
	"github.com/go-kit/kit/metrics"
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

func (s *instrumentingService) RecreateNode(ctx context.Context, nodeName, callerName string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RecreateNode").Add(1)
		s.requestLatency.With("method", "RecreateNode").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.RecreateNode(ctx, nodeName, callerName)
}
