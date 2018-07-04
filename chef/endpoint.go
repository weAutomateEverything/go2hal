package chef

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/telegram"
	"strconv"
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
		return nil, s.addEnvironmentToGroup(claim.RoomToken, req.EnvironmentName, req.FriendlyName)
	}
}

func makeGetEnvironmentForGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		return s.getEnvironmentForGroup(claim.RoomToken)
	}
}
func makeGetChefRecipesByGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req :=request.(string)
		i, _ := strconv.ParseUint(req, 10, 32)
		recipes,err:=s.getRecipesForGroup(uint32(i))
		l := make([]string, len(recipes))
		for x, i := range recipes {
			l[x] = i.FriendlyName
		}
		response = &recipeResponse{
			Recipes:l,
		}

		return response,nil

	}
}
func makeGetChefNodesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req :=request.(*chefNodeRequest)
		nodes:=s.FindNodesFromFriendlyNames(req.Recipe,req.Environment,req.Chat)
		res := make([]string, len(nodes))
		for i, x := range nodes {
			res[i] = x.Name
		}
		response = &nodeResponse{
			Nodes:res,
		}
		return response,nil

	}
}
func makeGetChefEnvironmentsByGroupEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req :=request.(string)
		groupid, _ := strconv.ParseUint(req, 10, 32)
		environments,err:=s.getEnvironmentForGroup(uint32(groupid))
		l := make([]string, len(environments))
		for x, i := range environments {
			l[x] = i.FriendlyName
		}
		response = &environmentResponse{
			Environments:l,
		}
		return response,nil
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
type recipeResponse struct{
	Recipes []string
}
type nodeResponse struct{
	Nodes []string
}
type environmentResponse struct{
	Environments []string
}
type chefNodeRequest struct{
	Recipe string
	Environment string
	Chat uint32
}