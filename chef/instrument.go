package chef

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

func (s *instrumentingService) sendDeliveryAlert(ctx context.Context, chatId uint32, message string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "sendDeliveryAlert").Add(1)
		s.requestLatency.With("method", "sendDeliveryAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.sendDeliveryAlert(ctx, chatId, message)
}
func (s *instrumentingService) FindNodesFromFriendlyNames(recipe, environment string) []Node {
	defer func(begin time.Time) {
		s.requestCount.With("method", "findNodesFromFriendlyNames").Add(1)
		s.requestLatency.With("method", "findNodesFromFriendlyNames").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.FindNodesFromFriendlyNames(recipe, environment)
}
