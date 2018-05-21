package alert

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

func makeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendAlert(ctx, ctx.Value("CHAT-ID").(uint32), req)
	}
}

func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToAlertGroup(ctx, ctx.Value("CHAT-ID").(uint32), req)
	}
}

func makeAlertErrorHandler(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendError(ctx, errors.New(req))

	}
}
