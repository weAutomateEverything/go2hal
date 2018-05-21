package halaws

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakeSendAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		s.SendAlert(ctx, ctx.Value("CHAT-ID").(uint32), req, "Manually Invoked", nil)
		return nil, nil
	}
}
