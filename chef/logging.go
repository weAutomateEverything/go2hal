package chef

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

func (s *loggingService)sendDeliveryAlert(message string){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendDeliveryAlert",
			"message",message,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.sendDeliveryAlert(message)

}
func (s *loggingService)findNodesFromFriendlyNames(recipe, environment string)[]node{
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "findNodesFromFriendlyNames",
			"recipe", recipe,
			"environment",environment,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.findNodesFromFriendlyNames(recipe,environment)

}
