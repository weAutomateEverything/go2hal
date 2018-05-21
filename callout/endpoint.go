package callout

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makeCalloutEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(SendCalloutRequest)
		if !ok {
			return nil, fmt.Errorf("request type not a SendCalloutRequest, received %v", request)
		}

		return nil, service.InvokeCallout(ctx, gokit.GetChatId(ctx), req.Title, req.Message, nil)
	}
}
