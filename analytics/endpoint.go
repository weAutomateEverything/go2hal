package analytics

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

func makeAnalyticsEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendAnalyticsAlert(req)
		return nil,nil
	}
}