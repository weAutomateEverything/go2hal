package alert

import (
	"fmt"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"time"
)

func NewGrpcService() AlertServer {
	return &grpcService{}
}

type grpcService struct {
	service  Service
	store    telegram.Store
	requests map[int64]Alert_RequestReplyServer
}

func (s *grpcService) RequestReply(a Alert_RequestReplyServer) error {
	for {
		in, err := a.Recv()
		if err != nil {
			s.service.SendError(a.Context(), fmt.Errorf("error receiving grpc alert: %v", err))
			continue
		}
		if in.Message != "" {
			_, err = s.service.SendAlertWithReply(a.Context(), uint32(in.GroupID), in.Message, in.CorrelationId)
			if err != nil {
				s.service.SendError(a.Context(), fmt.Errorf("error sending grpc alert message: %v", err))
				continue
			}
		}
		s.requests[in.GroupID] = a
	}
}

func (s grpcService) poll() {
	for {
		time.Sleep(2 * time.Second)
		replies, err := s.store.GetReplies()
		if err != nil {
			s.service.SendError(context.Background(), fmt.Errorf("error polling for replies %v", err))
		}
		for _, reply := range replies {
			groupId, err := s.store.GetUUID(reply.ChatId, "")
			if err != nil {
				s.service.SendError(context.Background(), fmt.Errorf("error looking up ground id for grpc alert service %v", err))
				continue
			}
			if s.requests[int64(groupId)] != nil {
				err = s.requests[int64(groupId)].Send(&AlertReply{
					CorrelationId: reply.CorrelationId,
					Message:       reply.Message,
				},
				)
			}
			if err != nil {
				s.service.SendError(context.Background(), fmt.Errorf("error sending grpc reply %v", err))
				continue
			}
			err = s.store.AcknowledgeReply(reply.ID)
			if err != nil {
				s.service.SendError(context.Background(), fmt.Errorf("error acknowledging reply %v", err))
				continue
			}
		}

	}
}
