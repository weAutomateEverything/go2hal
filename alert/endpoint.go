package alert

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/pkg/errors"
)


func makeAlertEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		err = s.SendAlert(req)
		return nil, err
	}
}


func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToAlertGroup(req)
	}
}

func makeHeartbeatMessageEncpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendHeartbeatGroupAlert(req)
	}
}

func makeImageHeartbeatEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToHeartbeatGroup(req)
	}
}

func makeBusinessAlertEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendNonTechnicalAlert(req)
		return nil, nil
	}
}

func makeAlertErrorHandler(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendError(errors.New(req))
		return nil,nil
	}
}

