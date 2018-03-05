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

func (s *loggingService) sendAppdynamicsAlert(ctx context.Context, message string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send appdynamics alert",
			"message", message,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.sendAppdynamicsAlert(ctx, message)
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
func (s *loggingService) addAppDynamicsQueue(name, application, metricPath string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "add appdynamics queue",
			"name", name,
			"application", application,
			"metricPath", metricPath,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.addAppDynamicsQueue(name, application, metricPath)
}
func (s *loggingService) executeCommandFromAppd(ctx context.Context, commandName, applicationID, nodeId string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "execute app dynamics command",
			"command_name", commandName,
			"applicatiom_id", applicationID,
			"node id", nodeId,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.executeCommandFromAppd(ctx, commandName, applicationID, nodeId)
}
