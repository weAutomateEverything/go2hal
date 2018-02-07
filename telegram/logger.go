package telegram

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) SendMessage(chatID int64, message string, messageID int) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendMessage",
			"chatId", chatID,
			"message", message,
			"messageId", messageID,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.SendMessage(chatID, message, messageID)
}
func (s *loggingService) SendMessagePlainText(chatID int64, message string, messageID int) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendMessagePlainText",
			"chatId", chatID,
			"message", message,
			"messageId", messageID,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.SendMessagePlainText(chatID, message, messageID)
}
func (s *loggingService) SendImageToGroup(image []byte, group int64) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendImageToGroup",
			"imageBytes", len(image),
			"groupId", group,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.SendImageToGroup(image, group)
}
func (s *loggingService) SendKeyboard(buttons []string, text string, chat int64) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "SendKeyboard",
			"buttons", buttons,
			"text", text,
			"chat", chat,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.SendKeyboard(buttons, text, chat)
}
func (s *loggingService) RegisterCommand(command Command) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RegisterCommand",
			"command", command,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RegisterCommand(command)
}
func (s *loggingService) RegisterCommandLet(commandlet Commandlet) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "RegisterCommandLet",
			"commandlet", commandlet,
			"took", time.Since(begin),
		)
	}(time.Now())
	s.Service.RegisterCommandLet(commandlet)
}
