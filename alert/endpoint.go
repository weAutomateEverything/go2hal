package alert

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/gokit"
)

func makeAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendAlert(ctx, gokit.GetChatId(ctx), req)
	}
}

func makeImageAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte)
		return nil, s.SendImageToAlertGroup(ctx, gokit.GetChatId(ctx), req)
	}
}

func makeAlertErrorHandler(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		return nil, s.SendError(ctx, errors.New(req))

	}
}

func makeDocumentAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(sendDocumentRequest)
		return nil, s.SendDocumentToAlertGroup(ctx, gokit.GetChatId(ctx), req.document, req.exension)
	}
}

func makeReplyAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SendReplyAlertMessageRequest)
		id, err := s.SendAlertWithReply(ctx, gokit.GetChatId(ctx), req.Message, req.CorrelationId)
		if err != nil {
			return nil, err
		}
		return SendReplyAlertMessageResponse{MessageId: id}, nil
	}
}

func makeGetRepliesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		replies, err := s.GetReplies(ctx, gokit.GetChatId(ctx))
		if err != nil {
			return nil, err
		}

		r := make([]Reply, len(replies))
		for i, reply := range replies {
			r[i] = Reply{
				Message:       reply.Message,
				CorrelationId: reply.CorrelationId,
				MessageId:     reply.ID.Hex(),
			}
		}
		return r, nil
	}
}

func makeAcknowledgeReplyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(acknowledgeReplyRequest)
		return nil, s.DeleteReply(ctx, gokit.GetChatId(ctx), req.messageId)
	}
}

type sendDocumentRequest struct {
	document []byte
	exension string
}

type acknowledgeReplyRequest struct {
	messageId string
}

//swagger:model
type SendReplyAlertMessageRequest struct {
	Message       string
	CorrelationId string
}

//swagger:model
type SendReplyAlertMessageResponse struct {
	MessageId int
}

//swager:model
type Reply struct {
	MessageId     string
	Message       string
	CorrelationId string
}
