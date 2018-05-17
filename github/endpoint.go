package github

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeSendAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendGithubMessage(ctx, req)
		return nil, nil
	}
}

