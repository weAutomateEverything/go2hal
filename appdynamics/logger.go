package appdynamics

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

func (s *loggingService) sendAppdynamicsAlert(ctx context.Context, chatId uint32, message string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send appdynamics alert",
			"message", message,
			"chat", chatId,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.sendAppdynamicsAlert(ctx, chatId, message)
}
func (s *loggingService) addAppdynamicsEndpoint(endpoint string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add appdynamics endpoint",
			"endpoint", endpoint,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.addAppdynamicsEndpoint(endpoint)
}
func (s *loggingService) addAppDynamicsQueue(chatId uint32, name, application, metricPath string) (err error) {
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
	return s.Service.addAppDynamicsQueue(chatId, name, application, metricPath)
}
func (s *loggingService) executeCommandFromAppd(ctx context.Context, chatId uint32, commandName, applicationID, nodeID string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "execute app dynamics command",
			"command_name", commandName,
			"applicatiom_id", applicationID,
			"node id", nodeID,
			"chat", chatId,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.executeCommandFromAppd(ctx, chatId, commandName, applicationID, nodeID)
}
