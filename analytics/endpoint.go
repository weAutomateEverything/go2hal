package analytics

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeAnalyticsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.SendAnalyticsAlert(ctx, req)
		return nil, nil
	}
}
