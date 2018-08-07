package alert

import (
	appd "appdynamics"
	"github.com/weAutomateEverything/go2hal/appdynamics/util"
	"golang.org/x/net/context"
)

func NewXray(s Service) Service {
	return &sXray{
		s,
	}
}

type sXray struct {
	Service
}

func (s sXray) SendAlert(ctx context.Context, chatId uint32, message string) (err error) {
	seg, ctx := util.Start("alert.SendAlert", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "chat", string(chatId))
	appd.AddUserDataToBT(seg, "message", message)
	err = s.Service.SendAlert(ctx, chatId, message)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return

}

func (s sXray) SendImageToAlertGroup(ctx context.Context, chatid uint32, image []byte) (err error) {
	seg, ctx := util.Start("alert.SendAlert", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "chatid", string(chatid))
	err = s.Service.SendImageToAlertGroup(ctx, chatid, image)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}

	return
}

func (s sXray) SendDocumentToAlertGroup(ctx context.Context, chatid uint32, document []byte, extension string) (err error) {
	seg, ctx := util.Start("alert.SendAlert", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "chat", string(chatid))
	appd.AddUserDataToBT(seg, "extension", extension)
	err = s.Service.SendDocumentToAlertGroup(ctx, chatid, document, extension)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return
}

func (s sXray) SendError(ctx context.Context, err error) (errout error) {
	seg, ctx := util.Start("alert.SendAlert", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	appd.AddUserDataToBT(seg, "err", err.Error())
	err = s.Service.SendError(ctx, err)

	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return
}

func (s sXray) SendErrorImage(ctx context.Context, image []byte) (err error) {
	seg, ctx := util.Start("alert.SendAlert", util.GetAppdUUID(ctx))
	defer appd.EndBT(seg)
	err = s.Service.SendErrorImage(ctx, image)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return
}
