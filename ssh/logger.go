package ssh

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

func (s *loggingService) ExecuteRemoteCommand(ctx context.Context, chatId uint32, commandName, address string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendSNMPMessage",
			"commandName", commandName,
			"address", address,
			"chat", chatId,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.ExecuteRemoteCommand(ctx, chatId, commandName, address)
}
