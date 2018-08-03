package alert

import (
	"context"
	"github.com/go-kit/kit/metrics"
	"strconv"
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

func (s *instrumentingService) SendAlert(ctx context.Context, chatId uint32, message string) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Add(1)
		s.requestLatency.With("method", "send_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Add(1)
		}
	}(time.Now())
	return s.Service.SendAlert(ctx, chatId, message)
}

func (s *instrumentingService) SendImageToAlertGroup(ctx context.Context, chatId uint32, image []byte) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_image_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Add(1)
		s.requestLatency.With("method", "send_image_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "send_image_alert", "chat", strconv.FormatUint(uint64(chatId), 10)).Add(1)
		}
	}(time.Now())
	return s.Service.SendImageToAlertGroup(ctx, chatId, image)
}

func (s *instrumentingService) SendError(ctx context.Context, err error) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send_error", "chat", "").Add(1)
		s.requestLatency.With("method", "send_error", "chat", "").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendError(ctx, err)
}
