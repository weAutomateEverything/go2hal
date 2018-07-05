package telegram

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"os"
	"time"
	"strconv"
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
	User int
	Chat int64
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
		id,err:=s.SendKeyboard(ctx,req.Options,req.Message,req.GroupId)
		return id,err
	}
}
func CustomClaimFactory() jwt.Claims {
	return &CustomClaims{}
}
func makeGetRoomEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req :=request.(string)
		i, _ := strconv.ParseUint(req, 10, 32)
		room,err:=s.GetRoomKey(uint32(i))
		response = &roomResponse{
			Id: room,
		}
		return response,nil

	}
}
type roomResponse struct {
	Id int64
}