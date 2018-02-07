package telegram

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

func (s *instrumentingService) SendMessage(chatID int64, message string, messageID int) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendMessage").Add(1)
		s.requestLatency.With("method", "SendMessage").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendMessage(chatID, message, messageID)
}
func (s *instrumentingService) SendMessagePlainText(chatID int64, message string, messageID int) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendMessagePlainText").Add(1)
		s.requestLatency.With("method", "SendMessagePlainText").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendMessagePlainText(chatID, message, messageID)
}
func (s *instrumentingService) SendImageToGroup(image []byte, group int64) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendImageToGroup").Add(1)
		s.requestLatency.With("method", "SendImageToGroup").Observe(time.Since(begin).Seconds())
	}(time.Now())
	return s.Service.SendImageToGroup(image, group)
}
func (s *instrumentingService) SendKeyboard(buttons []string, text string, chat int64) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "SendKeyboard").Add(1)
		s.requestLatency.With("method", "SendKeyboard").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.SendKeyboard(buttons, text, chat)
}
func (s *instrumentingService) RegisterCommand(command Command) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RegisterCommand").Add(1)
		s.requestLatency.With("method", "RegisterCommand").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RegisterCommand(command)
}
func (s *instrumentingService) RegisterCommandLet(commandlet Commandlet) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "RegisterCommandLet").Add(1)
		s.requestLatency.With("method", "RegisterCommandLet").Observe(time.Since(begin).Seconds())
	}(time.Now())
	s.Service.RegisterCommandLet(commandlet)
}
