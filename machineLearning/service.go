package machineLearning

import (
	"bytes"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"time"
)

type Service interface {
	StoreHTTPRequest(ctx context.Context, request *http.Request) context.Context
	StoreAction(ctx context.Context, action string, fields map[string]interface{})
}

type service struct {
	store Store
}

func NewServce(store Store) Service {
	return &service{store}
}

func (s *service) StoreHTTPRequest(ctx context.Context, request *http.Request) context.Context {
	body, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	id := s.store.SaveInputRecord("HTTP", time.Now(), map[string]interface{}{"body": string(body), "method": request.Method})
	return context.WithValue(ctx, "MSG_ID", id)
}

func (s *service) StoreAction(ctx context.Context, action string, fields map[string]interface{}) {
	val := ctx.Value("MSG_ID")
	if val != nil {
		s.store.SaveAction(val.(string), action, time.Now(), fields)
	}
}
