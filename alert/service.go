package alert

import (
	"errors"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strconv"
)

/*
NewService returns a new Alert Service
*/
func NewService(t telegram.Service, store telegram.Store) Service {

	return &service{
		store:    store,
		telegram: t,
	}
}

/*
Service interface
*/
type Service interface {
	SendAlert(ctx context.Context, chatId uint32, message string) error
	SendAlertWithReply(ctx context.Context, chatId uint32, message string, correlationId string) (int, error)
	SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) error
	SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) error

	SendError(ctx context.Context, err error) error
	SendErrorImage(ctx context.Context, image []byte) error

	GetReplies(ctx context.Context, chatId uint32) ([]telegram.Replies, error)
	DeleteReply(ctx context.Context, chatId uint32, messageId string) error
}

type service struct {
	telegram telegram.Service
	store    telegram.Store
}

func (s *service) GetReplies(ctx context.Context, chatId uint32) ([]telegram.Replies, error) {
	room, err := s.store.GetRoomKey(chatId)
	if err != nil {
		return nil, err
	}
	return s.store.GetReplies(room)

}

func (s *service) DeleteReply(ctx context.Context, chatId uint32, messageId string) error {
	return s.store.AcknowledgeReply(bson.ObjectIdHex(messageId))
}

func (s *service) SendAlertWithReply(ctx context.Context, chatId uint32, message string, correlationId string) (int, error) {
	group, err := s.store.GetRoomKey(chatId)
	if err != nil {
		return 0, err
	}
	msgId, err := s.telegram.SendMessageWithCorrelation(ctx, group, message, 0, correlationId)
	return msgId, err
}

//IMPL
func (s *service) SendAlert(ctx context.Context, chatid uint32, message string) error {
	group, err := s.store.GetRoomKey(chatid)
	if err != nil {
		return err
	}
	_, err = s.telegram.SendMessage(ctx, group, message, 0)
	return err
}

func (s *service) SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) error {

	group, err := s.store.GetRoomKey(chatid)
	if err != nil {
		return err
	}

	return s.telegram.SendImageToGroup(ctx, image, group)
}

func (s *service) SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) error {
	group, err := s.store.GetRoomKey(chatid)
	if err != nil {
		return err
	}

	return s.telegram.SendDocumentToGroup(ctx, document, extension, group)
}

func (s *service) SendError(ctx context.Context, err error) error {
	g := os.Getenv("ERROR_GROUP")
	group, errs := strconv.ParseInt(g, 10, 64)
	if errs != nil {
		err = errors.New("ERROR_GROUP has not been set to a valid chat group")
		return err
	}
	_, err = s.telegram.SendMessagePlainText(ctx, group, emoji.Sprintf(":poop: %s", err.Error()), 0)
	return err
}

func (s service) SendErrorImage(ctx context.Context, image []byte) error {
	g := os.Getenv("ERROR_GROUP")
	group, err := strconv.ParseUint(g, 10, 32)
	if err != nil {
		return err
	}
	return s.SendImageToAlertGroup(ctx, uint32(group), image)
}
