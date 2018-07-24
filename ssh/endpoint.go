package ssh

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makeAddCommandEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(addCommand)
		return nil, s.addCommand(gokit.GetChatId(ctx), r.Name, r.Command)
	}
}

func makeAddKeyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(addKey)
		return nil, s.addKey(gokit.GetChatId(ctx), r.UserName, r.Key)
	}
}

func makeExecuteCommandEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(executeCommand)
		return nil, s.ExecuteRemoteCommand(ctx, gokit.GetChatId(ctx), r.Command, r.Address)
	}
}

// swagger:model
type addCommand struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

// swagger:model
type addKey struct {
	UserName string `json:"user_name"`
	Key      string `json:"key"` //Base64 Encoded
}

// swagger:model
type executeCommand struct {
	Command string `json:"command"`
	Address string `json:"address"`
}
