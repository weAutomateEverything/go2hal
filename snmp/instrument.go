package snmp

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

func (s *instrumentingService) SendSNMPMessage(ctx context.Context, chat uint32) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendSNMPMessage").Add(1)
		s.requestLatency.With("method", "SendSNMPMessage").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.SendSNMPMessage(ctx, chat)
}
