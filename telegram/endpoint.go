package telegram

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
)

type CustomClaims struct {
	RoomToken uint32 `json:"roomToken"`
	jwt.StandardClaims
}

//swagger:model
type authRequestObject struct {
	RoomId uint32
	Name   string
}

//swagger:model
type authResponse struct {
	Key string
}

func makeTelegramAuthRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequestObject)
		token, err := s.requestAuthorisation(ctx, req.RoomId, req.Name)
		if err != nil {
			return
		}

		response = &authResponse{
			Key: token,
		}

		return
	}
}

func makeTelegramAuthPollEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		id := request.(string)
		room, err := s.pollAuthorisation(id)
		if err != nil {
			return
		}
		return makeToken(room)
	}
}

func CustomClaimFactory() jwt.Claims {
	return &CustomClaims{}
}
