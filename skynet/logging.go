package skynet

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

func (s loggingService) RecreateNode(ctx context.Context, nodeName, callerName string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendDeliveryAlert",
			"nodeName", nodeName,
			"callerName", callerName,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.RecreateNode(ctx, nodeName, callerName)
}
