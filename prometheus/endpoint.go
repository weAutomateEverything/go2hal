package prometheus

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makePrometheusAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		msg := request.(string)
		return nil, s.sendPrometheusAlert(ctx, gokit.GetChatId(ctx), msg)
	}
}
