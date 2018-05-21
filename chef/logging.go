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
func (s *loggingService) FindNodesFromFriendlyNames(recipe, environment string) []Node {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "findNodesFromFriendlyNames",
			"recipe", recipe,
			"environment", environment,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.FindNodesFromFriendlyNames(recipe, environment)

}
