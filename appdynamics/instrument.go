package appdynamics

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


func (s instrumentingService)sendAppdynamicsAlert(message string){
	defer func(begin time.Time) {
		s.requestCount.With("method", "sendAppdynamicsAlert").Add(1)
		s.requestLatency.With("method", "sendAppdynamicsAlert").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.sendAppdynamicsAlert(message)
}
func (s instrumentingService)addAppdynamicsEndpoint(endpoint string) error{
	defer func(begin time.Time) {
		s.requestCount.With("method", "addAppdynamicsEndpoint").Add(1)
		s.requestLatency.With("method", "addAppdynamicsEndpoint").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.addAppdynamicsEndpoint(endpoint)
}
func (s instrumentingService)addAppDynamicsQueue(name, application, metricPath string) error{
	defer func(begin time.Time) {
		s.requestCount.With("method", "addAppDynamicsQueue").Add(1)
		s.requestLatency.With("method", "addAppDynamicsQueue").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.addAppDynamicsQueue(name, application, metricPath)
}
func (s instrumentingService)executeCommandFromAppd(commandName, applicationID, nodeId string) error{
	defer func(begin time.Time) {
		s.requestCount.With("method", "executeCommandFromAppd").Add(1)
		s.requestLatency.With("method", "executeCommandFromAppd").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.executeCommandFromAppd(commandName, applicationID, nodeId)
}