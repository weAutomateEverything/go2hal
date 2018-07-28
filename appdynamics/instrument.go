package appdynamics

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

func (s instrumentingService) sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "sendAppdynamicsAlert").Add(1)
		s.requestLatency.With("method", "sendAppdynamicsAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.sendAppdynamicsAlert(ctx, chatId, message)
}
func (s instrumentingService) addAppdynamicsEndpoint(chat uint32, endpoint string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "addAppdynamicsEndpoint").Add(1)
		s.requestLatency.With("method", "addAppdynamicsEndpoint").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.addAppdynamicsEndpoint(chat, endpoint)
}
func (s instrumentingService) addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "addAppDynamicsQueue").Add(1)
		s.requestLatency.With("method", "addAppDynamicsQueue").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.addAppDynamicsQueue(ctx, chatId, name, application, metricPath)
}
func (s instrumentingService) executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "executeCommandFromAppd").Add(1)
		s.requestLatency.With("method", "executeCommandFromAppd").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.executeCommandFromAppd(ctx, chatId, commandName, applicationID, nodeID)
}
