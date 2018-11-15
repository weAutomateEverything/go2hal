package firstCall

import (
	"context"
	"errors"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/telegram"
)

func makeGetDefaultCalloutEndpoint(s CalloutFunction) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		_, number, err := s.GetFirstCallDetails(ctx, claim.RoomToken)
		if err != nil {
			return
		}
		response = DefaultCalloutRequest{
			PhoneNumber: number,
		}
		return
	}
}

func makeSetDefaultCalloutEndpoint(s DefaultCalloutService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		req, ok := request.(DefaultCalloutRequest)
		if !ok {
			err = errors.New("request is not oa Default Callout Request")
			return

		}
		err = s.setDefaultCallout(ctx, claim.RoomToken, req.PhoneNumber)
		return
	}
}

// swagger:model
type DefaultCalloutRequest struct {
	PhoneNumber string
}
