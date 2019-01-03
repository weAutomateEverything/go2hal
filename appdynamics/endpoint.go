package appdynamics

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

// Request object to add a new IBM MQ Queue to be monitored from App Dynamics
// swagger:model
type AddAppdynamicsQueueEndpointRequest struct {
	Name         string
	Application  string
	Metricpath   string
	IgnorePrefix []string `json:"ignore_prefix"`
}

//swagger:model
type AddAppdynamicsEndpointRequest struct {
	Endpoint string
}

// swagger:model
type ExecuteAppDynamicsCommandRequest struct {
	CommandName   string `json:"command_name"`
	NodeID        string `json:"node_id"`
	ApplicationID string `json:"application_id"`
}

type BusinessAlertRequest struct {
	Severity, Type, DisplayName, SummaryMessage string
}

func makeAppDynamicsAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AppdynamicsMessage)
		return nil, s.sendAppdynamicsAlert(ctx, gokit.GetChatId(ctx), req)
	}
}

func makeAddAppdynamicsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddAppdynamicsEndpointRequest)
		return nil, s.addAppdynamicsEndpoint(gokit.GetChatId(ctx), req.Endpoint)

	}
}

func makExecuteCommandFromAppdynamics(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ExecuteAppDynamicsCommandRequest)
		return nil, s.executeCommandFromAppd(ctx, gokit.GetChatId(ctx), req.CommandName, req.ApplicationID, req.NodeID)

	}
}
