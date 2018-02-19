package alert

import (
	"github.com/go-kit/kit/metrics"
	"time"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	errorCount     metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func NewInstrumentService(counter metrics.Counter, errorCount metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		errorCount:     errorCount,
		Service:        s,
	}
}

func (s *instrumentingService) SendAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_alert").Add(1)
		s.requestLatency.With("method", "send_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlert(message)
}

func (s *instrumentingService) SendNonTechnicalAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_technical_alert").Add(1)
		s.requestLatency.With("method", "send_technical_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_technical_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendNonTechnicalAlert(message)
}

func (s *instrumentingService) SendHeartbeatGroupAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_heartbeat_alert").Add(1)
		s.requestLatency.With("method", "send_heartbeat_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_heartbeat_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendHeartbeatGroupAlert(message)
}

func (s *instrumentingService) SendImageToAlertGroup(image []byte) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_image_alert").Add(1)
		s.requestLatency.With("method", "send_image_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_image_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendImageToAlertGroup(image)
}

func (s *instrumentingService) SendImageToHeartbeatGroup(image []byte) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_image_heartbeat").Add(1)
		s.requestLatency.With("method", "send_image_heartbeat").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_image_heartbeat").Add(1)
		}
	}(time.Now())
	return s.Service.SendImageToHeartbeatGroup(image)
}

func (s *instrumentingService) SendError(err error) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_error").Add(1)
		s.requestLatency.With("method", "send_error").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendError(err)
}
