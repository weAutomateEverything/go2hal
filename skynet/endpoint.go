package skynet

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type SkynetRebuildRequest struct {
	NodeName string `json:"Nodename"`
	User     string `json:"User"`
}


func makeSkynetRebuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SkynetRebuildRequest)
		s.RecreateNode(req.NodeName,req.User)
		return nil,err
	}
}