package chef

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type AddChefClientRequest struct {
	Name, Key, URL string
}

func makeChefDeliveryAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendDeliveryAlert(req)
		return nil, nil
	}
}
