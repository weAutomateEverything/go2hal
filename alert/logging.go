package alert

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

func (s *loggingService) SendAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlert(message)

}
func (s *loggingService) SendNonTechnicalAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send non technical alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendNonTechnicalAlert(message)

}
func (s *loggingService) SendHeartbeatGroupAlert(message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send heartbeat group alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendHeartbeatGroupAlert(message)

}
func (s *loggingService) SendImageToAlertGroup(image []byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send image to alert group",
			"image length", len(image),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendImageToAlertGroup(image)

}
func (s *loggingService) SendImageToHeartbeatGroup(image []byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send image to heartbeat group",
			"image length", len(image),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendImageToHeartbeatGroup(image)

}
func (s *loggingService) SendError(err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send error",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	s.Service.SendError(err)

}
