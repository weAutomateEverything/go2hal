package prometheus

import (
	"github.com/go-kit/kit/log"
	"time"
)

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) sendPrometheusAlert(chat uint32, body string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendPrometheusAlert",
			"body", body,
			"chat", chat,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.sendPrometheusAlert(chat, body)

}
