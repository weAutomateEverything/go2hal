package alert

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)
//request for keyboard alerts.
type KeyboardAlertRequest struct {
	Nodes   []string `json:"Nodes"`
}

func makeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendAlert(ctx, req)
	}
}

func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToAlertGroup(ctx, req)
	}
}

func makeHeartbeatMessageEncpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendHeartbeatGroupAlert(ctx, req)
	}
}

func makeImageHeartbeatEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToHeartbeatGroup(ctx, req)
	}
}

func makeBusinessAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendNonTechnicalAlert(ctx, req)
	}
}

func makeAlertErrorHandler(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendError(ctx, errors.New(req))

	}

}

func makeKeyboardRecipeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(KeyboardAlertRequest)
		return nil, s.SendAlertKeyboardRecipe(ctx,req.Nodes)
	}
}
func makeEnvironmentAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(KeyboardAlertRequest)
		return nil, s.SendAlertEnvironment(ctx, req.Nodes)
	}
}
func makeNodesAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(KeyboardAlertRequest)
		return nil, s.SendAlertNodes(ctx, req.Nodes)
	}
}