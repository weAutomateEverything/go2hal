package appdynamics

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type mqLoggingService struct {
	logger log.Logger
	MqService
}

func NewMqLoggingService(logger log.Logger, s MqService) MqService {
	return &mqLoggingService{logger, s}
}

func (s *mqLoggingService) addAppDynamicsQueue(ctx context.Context, chatId uint32, name, application, metricPath string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add appdynamics queue",
			"name", name,
			"application", application,
			"metricPath", metricPath,
			"error", err,
			"chat", chatId,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.MqService.addAppDynamicsQueue(ctx, chatId, name, application, metricPath)
}
