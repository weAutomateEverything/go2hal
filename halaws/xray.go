package halaws

import (
	appd "appdynamics"
	"github.com/weAutomateEverything/go2hal/appdynamics/util"
	"golang.org/x/net/context"
)

func NewXray(service Service) Service {
	return &x{
		service: service,
	}
}

type x struct {
	service Service
}

func (s *x) SendAlert(ctx context.Context, chatId uint32, destination string, name string, variables map[string]string) (err error) {
	handler, ctx := util.Start("halaws.SendAlert", "")
	defer appd.EndBT(handler)
	appd.AddUserDataToBT(handler, "chatid", string(chatId))
	appd.AddUserDataToBT(handler, "destination", destination)
	appd.AddUserDataToBT(handler, "name", name)
	err = s.service.SendAlert(ctx, chatId, destination, name, variables)
	if err != nil {
		util.AddErrorToAppDynamics(ctx, err)
	}
	return
}
