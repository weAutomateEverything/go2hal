package sensu

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

func (s *instrumentingService) handleSensu(ctx context.Context, chatId uint32, sensu SensuMessageRequest) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "handleSensu").Add(1)
		s.requestLatency.With("method", "handleSensu").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.handleSensu(ctx, chatId, sensu)
}
