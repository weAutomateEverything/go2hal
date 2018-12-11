package ssh

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/telegram"
)

func makeAddCommandEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		r := request.(addCommand)
		return nil, s.addCommand(claim.RoomToken, r.Name, r.Command)
	}
}

func makeAddKeyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)

		r := request.(addKey)
		return nil, s.addKey(claim.RoomToken, r.UserName, r.Key)
	}
}

func makeAddServerEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)

		r := request.(addServer)
		return nil, s.addServer(claim.RoomToken, r.Address, r.Description)
	}
}

func makeExecuteCommandEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)

		r := request.(executeCommand)
		return nil, s.ExecuteRemoteCommand(ctx, claim.RoomToken, r.Command, r.Address)
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

//swagger:model
type addServer struct {
	Address     string `json:"address"`
	Description string `json:"description"`
}
