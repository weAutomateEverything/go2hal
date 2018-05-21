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

func (s *loggingService) SendAlert(ctx context.Context, chatId uint32, message string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send alert",
			"message", message,
			"chat", chatId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendAlert(ctx, chatId, message)

}

func (s *loggingService) SendImageToAlertGroup(ctx context.Context, chatId uint32, image []byte) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send image to alert group",
			"image_length", len(image),
			"chat", chatId,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendImageToAlertGroup(ctx, chatId, image)

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
