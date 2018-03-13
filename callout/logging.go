package callout

import (
	"context"
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

func (s loggingService) InvokeCallout(ctx context.Context, title, message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "InvokeCallout",
			"title", title,
			"message", message,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.InvokeCallout(ctx, title, message)
}

func (s loggingService) getFirstCallName(ctx context.Context) (name string, phone string, error error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "getFirstCallName",
			"name", name,
			"phone", phone,
			"error", error,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.getFirstCall(ctx)
}
