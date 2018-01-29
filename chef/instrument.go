package chef

import (
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

func (s *instrumentingService)sendDeliveryAlert(message string){
	defer func(begin time.Time) {
		s.requestCount.With("method", "sendDeliveryAlert").Add(1)
		s.requestLatency.With("method", "sendDeliveryAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.sendDeliveryAlert(message)
}
func (s *instrumentingService)findNodesFromFriendlyNames(recipe, environment string)[]node{
	defer func(begin time.Time) {
		s.requestCount.With("method", "findNodesFromFriendlyNames").Add(1)
		s.requestLatency.With("method", "findNodesFromFriendlyNames").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.findNodesFromFriendlyNames(recipe,environment)
}

