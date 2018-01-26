package alert

import (
	"github.com/zamedic/go2hal/telegram"
	"gopkg.in/kyokomi/emoji.v1"
	"log"
)

type Service interface {
	SendAlert(message string) error
	SendNonTechnicalAlert(message string) error
	SendHeartbeatGroupAlert(message string) error
	SendImageToAlertGroup(image []byte) error
	SendImageToHeartbeatGroup(image []byte) error
	SendError(err error)
}

type service struct {
	telegram telegram.Service
	store Store
}

func NewService(t telegram.Service) Service {


	return &service{
		telegram: t,
	}
}

//IMPL

func (s *service) SendAlert(message string) error {
	alertGroup, err := s.store.alertGroup()
	if err != nil {
		return err
	}
	err = s.telegram.SendMessage(alertGroup, message, 0)
	return err
}

func (s *service) SendNonTechnicalAlert(message string) error {
	return nil
}

func (s *service) SendImageToAlertGroup(image []byte) error {

	alertGroup, err := s.store.alertGroup()
	if err != nil {
		s.SendError(err)
		return err
	}

	return s.telegram.SendImageToGroup(image, alertGroup)
}

func (s *service) SendImageToHeartbeatGroup(image []byte) error {
	group, err := s.store.heartbeatGroup()
	if err != nil {
		s.SendError(err)
		return err
	}

	return s.telegram.SendImageToGroup(image, group)

}

func (s service) SendHeartbeatGroupAlert(message string) error {
	group, err := s.store.heartbeatGroup()
	if err != nil {
		s.SendError(err)
		return err
	}

	return s.telegram.SendMessage(group, message, 0)
}

func (s *service) SendError(err error) {
	log.Println(err.Error())
	s.SendHeartbeatGroupAlert(emoji.Sprintf(":poop: %s", err.Error()))
}
