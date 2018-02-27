package user

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

func (s *loggingService) parseInputRequest(in string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "parseInputRequest",
			"Input", in,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.parseInputRequest(in)
}
