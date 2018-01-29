package snmp

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

func (s *instrumentingService)	SendSNMPMessage(){
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendSNMPMessage").Add(1)
		s.requestLatency.With("method", "SendSNMPMessage").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.SendSNMPMessage()
}

