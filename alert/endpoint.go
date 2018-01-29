package alert

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"encoding/base64"
	"errors"
)

func makeAlertEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		err = s.SendAlert(req)
		return nil, err
	}
}

type imageAlertMessage struct {
	Message, Image string
	InternalError bool
}

func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(imageAlertMessage)
		if req.Image != "" {
			b, err := base64.StdEncoding.DecodeString(req.Image)
			if err != nil {
				return nil, nil
			}
			if req.InternalError {
				err = s.SendImageToHeartbeatGroup(b)
			} else {
				err = s.SendImageToAlertGroup(b)
			}
		}
		if req.InternalError{
			s.SendError(errors.New(req.Message))
		} else {
			s.SendAlert(req.Message)
		}
		return nil, nil
	}
}

func makeBusinessAlertEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendNonTechnicalAlert(req)
		return nil, nil
	}
}
