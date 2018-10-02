package appdynamics

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"time"
)

type mqInstrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	mqService      MqService
}

func NewMqInstrumentService(counter metrics.Counter, latency metrics.Histogram, s MqService) MqService {
	return &mqInstrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		mqService:      s,
	}
}

func (s mqInstrumentingService) addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string, ignorePrefix []string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "addAppDynamicsQueue").Add(1)
		s.requestLatency.With("method", "addAppDynamicsQueue").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.mqService.addAppDynamicsQueue(ctx, chatId, name, application, metricPath, ignorePrefix)
}
