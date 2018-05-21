package github

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeSendAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendGithubMessage(ctx, ctx.Value("CHAT-ID").(uint32), req)
		return nil, nil
	}
}
