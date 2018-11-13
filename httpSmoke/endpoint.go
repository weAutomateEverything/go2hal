package httpSmoke

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/telegram"
)

func addHTTPEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AddHttpRequest)
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)

		var v []parameters
		for _, val := range req.Parameters {
			v = append(v, parameters{
				Name:  val.Name,
				Value: val.Value,
			})
		}

		return nil, s.addHttpEndpoint(ctx, req.Name, req.URL, req.Method, v, req.Threshold, claim.RoomToken)
	}
}

func getHTTPForGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		return s.getEndpoints(claim.RoomToken)
	}
}

func deleteHTTPEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

// swagger:model
type AddHttpRequest struct {
	Name       string
	Method     string
	URL        string
	Threshold  int
	Parameters []struct {
		Name  string
		Value string
	}
}
