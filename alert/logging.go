package alert

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

func (s *loggingService) SendAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlert(ctx, message)

}
func (s *loggingService) SendNonTechnicalAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send non technical alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendNonTechnicalAlert(ctx, message)

}
func (s *loggingService) SendHeartbeatGroupAlert(ctx context.Context, message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send heartbeat group alert",
			"message", message,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendHeartbeatGroupAlert(ctx, message)

}
func (s *loggingService) SendImageToAlertGroup(ctx context.Context, image []byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send image to alert group",
			"image_length", len(image),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendImageToAlertGroup(ctx, image)

}
func (s *loggingService) SendImageToHeartbeatGroup(ctx context.Context, image []byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send image to heartbeat group",
			"image_length", len(image),
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendImageToHeartbeatGroup(ctx, image)

}
func (s *loggingService) SendError(ctx context.Context, err error) (e error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send error",
			"took", time.Since(begin),
			"send_err", err,
			"response_err", e,
		)
	}(time.Now())
	return s.Service.SendError(ctx, err)

}
func (s *loggingService) SendAlertKeyboardRecipe(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send keyboard recipe",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlertKeyboardRecipe(ctx, nodes)

}
func (s *loggingService) SendAlertEnvironment(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send keyboard environment",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlertEnvironment(ctx, nodes)

}
func (s *loggingService) SendAlertNodes(ctx context.Context, nodes []string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send keyboard nodes",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlertNodes(ctx, nodes)

}
