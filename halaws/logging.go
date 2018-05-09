package halaws

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) SendAlert(ctx context.Context, destination string, name string, variables map[string]string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "InvokeCallout",
			"destination", destination,
			"error", err,
			"variables", variables,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.SendAlert(ctx, destination, name, variables)
}
