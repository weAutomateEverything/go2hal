package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeBulkUserUploadEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.parseInputRequest(req)
	}
}
