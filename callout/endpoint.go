package callout

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makeCalloutEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*SendCalloutRequest)
		if !ok {
			return nil, fmt.Errorf("request type not a SendCalloutRequest, received %v", request)
		}

		return nil, service.InvokeCallout(ctx, gokit.GetChatId(ctx), req.Title, req.Message)
	}
}

// Callout Request
//
// swagger:model
type SendCalloutRequest struct {
	// Title for the JIRA ticket
	//
	// required: true
	Title string `json:"title"`

	// Message that will be used for the Telegram Alert, the Jira Ticket and the Alexa Callout
	//
	// required: true
	Message string `json:"message"`
}
