package alert

import (
	"context"
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

func (s *instrumentingService) SendAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_alert").Add(1)
		s.requestLatency.With("method", "send_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlert(ctx, message)
}

func (s *instrumentingService) SendNonTechnicalAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_technical_alert").Add(1)
		s.requestLatency.With("method", "send_technical_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_technical_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendNonTechnicalAlert(ctx, message)
}

func (s *instrumentingService) SendHeartbeatGroupAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_heartbeat_alert").Add(1)
		s.requestLatency.With("method", "send_heartbeat_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_heartbeat_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendHeartbeatGroupAlert(ctx, message)
}

func (s *instrumentingService) SendImageToAlertGroup(ctx context.Context, image []byte) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_image_alert").Add(1)
		s.requestLatency.With("method", "send_image_alert").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_image_alert").Add(1)
		}
	}(time.Now())
	return s.Service.SendImageToAlertGroup(ctx, image)
}

func (s *instrumentingService) SendImageToHeartbeatGroup(ctx context.Context, image []byte) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_image_heartbeat").Add(1)
		s.requestLatency.With("method", "send_image_heartbeat").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_image_heartbeat").Add(1)
		}
	}(time.Now())
	return s.Service.SendImageToHeartbeatGroup(ctx, image)
}

func (s *instrumentingService) SendError(ctx context.Context, err error) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_error").Add(1)
		s.requestLatency.With("method", "send_error").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendError(ctx, err)
}
func (s *instrumentingService) SendAlertKeyboardRecipe(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_keyboard_recipe").Add(1)
		s.requestLatency.With("method", "send_keyboard_recipe").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_keyboard_recipe").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlertKeyboardRecipe(ctx, nodes)
}
func (s *instrumentingService) SendAlertEnvironment(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_keyboard_environment").Add(1)
		s.requestLatency.With("method", "send_keyboard_environment").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_keyboard_environment").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlertEnvironment(ctx, nodes)
}
func (s *instrumentingService) SendAlertNodes(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_keyboard_nodes").Add(1)
		s.requestLatency.With("method", "send_keyboard_nodes").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_keyboard_nodes").Add(1)
		}
	}(time.Now())
	return s.Service.SendAlertNodes(ctx, nodes)
}
