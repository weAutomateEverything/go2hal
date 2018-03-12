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

func (s loggingService) SendAlert(ctx context.Context, destination string, name string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "InvokeCallout",
			"destination", destination,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.SendAlert(ctx, destination, name)
}
