package chef

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

func (s *loggingService) sendDeliveryAlert(ctx context.Context, chatId uint32, message string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendDeliveryAlert",
			"message", message,
			"chat", chatId,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.sendDeliveryAlert(ctx, chatId, message)

}
func (s *loggingService) FindNodesFromFriendlyNames(ctx context.Context, recipe, environment string, chat uint32) []Node {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "findNodesFromFriendlyNames",
			"recipe", recipe,
			"chat", chat,
			"environment", environment,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.FindNodesFromFriendlyNames(ctx, recipe, environment, chat)

}
