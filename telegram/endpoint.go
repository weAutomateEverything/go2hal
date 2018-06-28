package telegram

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"os"
	"time"
	"strconv"
	"log"
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
	GroupId int64
}
type setStateRequest struct{
	User string
	Chat string
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
		user, err := strconv.Atoi(req.User)
		chat, err := strconv.ParseInt(req.Chat, 10, 32)
		err=s.SetState(user,chat,req.State,req.Field)
		if(err!=nil){
			return nil,err
		}
		return "success",nil
	}
}
func makeSendKeyboardEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(sendKeyBoardRequest)
		log.Println(req)
		id,err:=s.SendKeyboard(ctx,req.Options,req.Message,req.GroupId)
		return id,err
	}
}
func CustomClaimFactory() jwt.Claims {
	return &CustomClaims{}
}
