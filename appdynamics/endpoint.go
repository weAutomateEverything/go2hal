package appdynamics

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

type AddAppdynamicsQueueEndpointRequest struct {
	Name        string
	Application string
	Metricpath  string
}

type AddAppdynamicsEndpointRequest struct {
	Endpoint string
}

// swagger:model
type ExecuteAppDynamicsCommandRequest struct {
	CommandName   string `json:"command_name"`
	NodeID        string `json:"node_id"`
	ApplicationID string `json:"application_id"`
}

// swagger:model
type appdynamicsMessage struct {
	Environment string `json:"environment"`
	Policy      struct {
		TriggerTime string `json:"triggerTime"`
		Name        string `json:"name"`
	}
	Events []struct {
		Severity    string `json:"severity"`
		Application struct {
			Name string `json:"name"`
		}
		Tier struct {
			Name string `json:"name"`
		}
		Node struct {
			Name string `json:"name"`
		}
		DisplayName  string `json:"displayName"`
		EventMessage string `json:"eventMessage"`
	}
}

type BusinessAlertRequest struct {
	Severity, Type, DisplayName, SummaryMessage string
}

func makeAppDynamicsAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendAppdynamicsAlert(ctx, gokit.GetChatId(ctx), req)
		return nil, nil
	}
}

func makeAddAppdynamicsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddAppdynamicsEndpointRequest)
		return nil, s.addAppdynamicsEndpoint(gokit.GetChatId(ctx), req.Endpoint)

	}
}
func makeAddAppdynamicsQueueEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddAppdynamicsQueueEndpointRequest)
		return nil, s.addAppDynamicsQueue(ctx, gokit.GetChatId(ctx), req.Name, req.Application, req.Metricpath)
	}
}
func makExecuteCommandFromAppdynamics(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ExecuteAppDynamicsCommandRequest)
		return nil, s.executeCommandFromAppd(ctx, gokit.GetChatId(ctx), req.CommandName, req.ApplicationID, req.NodeID)

	}
}
