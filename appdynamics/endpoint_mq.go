package appdynamics

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
	"golang.org/x/net/context"
)

func makeAddAppdynamicsQueueEndpoint(s MqService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddAppdynamicsQueueEndpointRequest)
		return nil, s.addAppDynamicsQueue(ctx, gokit.GetChatId(ctx), req.Name, req.Application, req.Metricpath)
	}
}
