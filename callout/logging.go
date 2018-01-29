package callout

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

func (s loggingService)	InvokeCallout(title, message string){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "InvokeCallout",
			"title", title,
			"message",message,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.InvokeCallout(title,message)
}

func (s loggingService)	getFirstCallName() (name string, error error){
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "getFirstCallName",
			"name", name,
			"error",error,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.getFirstCallName()
}