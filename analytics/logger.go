package analytics

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) SendAnalyticsAlert(message string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send non technical alert",
			"message", message,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.SendAnalyticsAlert(message)
}
