package machineLearning

import (
	"bytes"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"time"
)

/*
Service forms the basis to build a machine learning data set.

The idea is to log all automated requests into the system and correlate them with actions taken by HAL.
Doing so, we will have a link between request and response for audit purposes, however the plan is to use the data
to try and implement some machine learning, where we can teach HAL what actions to take based on input.
*/
type Service interface {
	StoreHTTPRequest(ctx context.Context, request *http.Request) context.Context
	StoreAction(ctx context.Context, action string, fields map[string]interface{})
}

type service struct {
	store Store
}

type key int

const (
	keyPrincipalID key = iota
)

/*
NewService returns a new Machine Learning Service.
*/
func NewService(store Store) Service {
	return &service{store}
}

func (s *service) StoreHTTPRequest(ctx context.Context, request *http.Request) context.Context {
	body, _ := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	id := s.store.SaveInputRecord("HTTP", time.Now(), map[string]interface{}{"body": string(body), "method": request.Method})
	return context.WithValue(ctx, keyPrincipalID, id)
}

func (s *service) StoreAction(ctx context.Context, action string, fields map[string]interface{}) {
	val := ctx.Value(keyPrincipalID)
	if val != nil {
		s.store.SaveAction(val.(string), action, time.Now(), fields)
	}
}
