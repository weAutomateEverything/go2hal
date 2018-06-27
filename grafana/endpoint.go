package grafana

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makeGrafanaAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.sendGrafanaAlert(gokit.GetChatId(ctx), req)
	}
}
