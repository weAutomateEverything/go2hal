package telegram

import (
	appd "appdynamics"
	"context"
	"github.com/weAutomateEverything/go2hal/appdynamics/util"
)

func NewXray(s Service) Service {
	return sXray{
		s,
	}
}

type sXray struct {
	Service
}

func (s sXray) SendMessage(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	seg, ctx := util.Start("telegram.SendMessage", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "chat", string(chatID))
	appd.AddUserDataToBT(seg, "message", message)
	appd.AddUserDataToBT(seg, "messageid", string(messageID))
	msgid, err = s.Service.SendMessage(ctx, chatID, message, messageID)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return

}
func (s sXray) SendMessagePlainText(ctx context.Context, chatID int64, message string, messageID int) (msgid int, err error) {
	seg, ctx := util.Start("telegram.SendMessagePlainText", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "chatid", string(chatID))
	appd.AddUserDataToBT(seg, "message", string(message))
	appd.AddUserDataToBT(seg, "messageid", string(messageID))
	msgid, err = s.Service.SendMessagePlainText(ctx, chatID, message, messageID)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return
}
func (s sXray) SendImageToGroup(ctx context.Context, image []byte, group int64) (err error) {
	seg, ctx := util.Start("telegram.SendMessagePlainText", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "group", string(group))
	err = s.Service.SendImageToGroup(ctx, image, group)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return

}
func (s sXray) SendDocumentToGroup(ctx context.Context, document []byte, extension string, group int64) (err error) {
	seg, ctx := util.Start("telegram.SendMessagePlainText", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "extension", extension)
	appd.AddUserDataToBT(seg, "group", string(group))
	err = s.Service.SendDocumentToGroup(ctx, document, extension, group)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return

}

func (s sXray) SendKeyboard(ctx context.Context, buttons []string, text string, chat int64) (message int, err error) {
	seg, ctx := util.Start("telegram.SendMessagePlainText", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "text", string(text))
	appd.AddUserDataToBT(seg, "chat", string(chat))
	message, err = s.Service.SendKeyboard(ctx, buttons, text, chat)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}

	return

}
func (s sXray) RegisterCommand(command Command) {
	s.Service.RegisterCommand(command)
}
func (s sXray) RegisterCommandLet(commandlet Commandlet) {
	s.Service.RegisterCommandLet(commandlet)
}

func (s sXray) requestAuthorisation(ctx context.Context, chat uint32, name string) (string, error) {
	return s.Service.requestAuthorisation(ctx, chat, name)
}
func (s sXray) pollAuthorisation(token string) (uint32, error) {
	return s.Service.pollAuthorisation(token)
}
