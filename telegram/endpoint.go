package telegram

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"os"
	"time"
)

type CustomClaims struct {
	RoomToken uint32 `json:"roomToken"`
	jwt.StandardClaims
}

type authRequestObject struct {
	RoomId uint32
	Name   string
}

type authResponse struct {
	Key string
}
type sendKeyBoardRequest struct{
	Options []string
	Message string
	GroupId uint32
}
type setStateRequest struct{
	User int
	Chat uint32
	State string
	Field []string
}
func makeTelegramAuthRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authRequestObject)
		token, err := s.requestAuthorisation(req.RoomId, req.Name)
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
		claims := CustomClaims{
			room,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 120).Unix(),
				IssuedAt:  jwt.TimeFunc().Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString([]byte(os.Getenv("JWT_KEY")))
	}
}
func makeSetStateRequestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(setStateRequest)
		err=s.SetState(req.User,req.Chat,req.State,req.Field)
		if(err!=nil){
			return nil,err
		}
		return nil,nil
	}
}
func makeSendKeyboardEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(sendKeyBoardRequest)
		err=s.SendKeyboardGroup(ctx,req.Options,req.Message,req.GroupId)
		return nil,err
	}
}
func CustomClaimFactory() jwt.Claims {
	return &CustomClaims{}
}


