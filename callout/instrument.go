package callout

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

func (s instrumentingService) InvokeCallout(title, message string) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "InvokeCallout").Add(1)
		s.requestLatency.With("method", "InvokeCallout").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.InvokeCallout(title, message)
}

func (s instrumentingService) getFirstCallName() (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "getFirstCallName").Add(1)
		s.requestLatency.With("method", "getFirstCallName").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.getFirstCallName()
}
