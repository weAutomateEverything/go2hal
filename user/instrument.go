package user

import (
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	errorCount     metrics.Counter
	Service
}

/*
NewInstrumentService creates a prometheus service that will log request count, error count and latency
*/
func NewInstrumentService(counter metrics.Counter, errorCount metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		errorCount:     errorCount,
		Service:        s,
	}
}

func (s *instrumentingService) parseInputRequest(in string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "parseInputRequest").Add(1)
		s.requestLatency.With("method", "parseInputRequest").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "parseInputRequest").Add(1)
		}
	}(time.Now())
	return s.Service.parseInputRequest(in)
}
