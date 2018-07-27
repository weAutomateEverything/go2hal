package grafana

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"time"
)

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

type loggingService struct {
	logger log.Logger
	Service
}

func (s *loggingService) sendGrafanaAlert(ctx context.Context, chat uint32, body string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendGrafanaAlert",
			"body", body,
			"chat", chat,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.sendGrafanaAlert(ctx, chat, body)

}
