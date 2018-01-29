package sensu

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/go-kit/kit/log"
)

type SensuMessageRequest struct {
	Text        string            `json:"text"`
	IconURL     string            `json:"icon_url"`
	Attachments []sensuAttachment `json:"attachments"`
}

type sensuAttachment struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}



func makeSensuEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SensuMessageRequest)
		s.handleSensu(req)
		return nil,nil
	}
}

func loggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("request", request)
			defer logger.Log("msg", "called endpoint")
			return next(ctx, request)
		}
	}
}


