package halaws

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

func MakeSendAlertEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendAlert(req)
		return nil,nil
	}
}
