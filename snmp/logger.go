package snmp

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

func (s *loggingService) SendSNMPMessage(ctx context.Context) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendSNMPMessage",
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.SendSNMPMessage(ctx)

}
