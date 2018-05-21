package halaws

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func MakeSendAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)

		s.SendAlert(ctx, gokit.GetChatId(ctx), req, "Manually Invoked", nil)
		return nil, nil
	}
}
