package alert

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

func makeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendAlert(req)
	}
}

func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToAlertGroup(req)
	}
}

func makeHeartbeatMessageEncpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendHeartbeatGroupAlert(req)
	}
}

func makeImageHeartbeatEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToHeartbeatGroup(req)
	}
}

func makeBusinessAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendNonTechnicalAlert(req)
	}
}

func makeAlertErrorHandler(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendError(errors.New(req))

	}
}
