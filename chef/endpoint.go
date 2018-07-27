package chef

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/telegram"
)

type AddChefClientRequest struct {
	Name, Key, URL string
}

func makeChefDeliveryAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendDeliveryAlert(ctx, gokit.GetChatId(ctx), req)
		return nil, nil
	}
}

func makeAddRecipeToGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(*addRecipeRequest)
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		return nil, s.addRecipeToGroup(ctx, claim.RoomToken, r.RecipeName, r.FriendlyName)
	}
}

func makeGetAllGrouRecipesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		return s.getRecipesForGroup(claim.RoomToken)
	}
}

func makeAddEnvironmentToGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		req := request.(*addEnvironmentRequest)
		return nil, s.addEnvironmentToGroup(ctx, claim.RoomToken, req.EnvironmentName, req.FriendlyName)
	}
}

func makeGetEnvironmentForGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		return s.getEnvironmentForGroup(claim.RoomToken)
	}
}

type addRecipeRequest struct {
	RecipeName   string
	FriendlyName string
}

type addEnvironmentRequest struct {
	EnvironmentName string
	FriendlyName    string
}
