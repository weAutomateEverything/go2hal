package jira

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) CreateJira(ctx context.Context, title, description string, name string) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "CreateJira",
			"title", title,
			"description", description,
			"name", name,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.CreateJira(ctx, title, description, name)
}
