package ssh

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

func (s *instrumentingService) parseInputRequest(ctx context.Context, chatId uint32, commandName, address string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "parseInputRequest").Add(1)
		s.requestLatency.With("method", "parseInputRequest").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.ExecuteRemoteCommand(ctx, chatId, commandName, address)
}
