package util

import (
	appd "appdynamics"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"log"
)

func Start(function, correlation string) (appd.BtHandle, context.Context) {
	bt := appd.StartBT(function, correlation)
	return bt, GetContext(bt)
}

func GetContext(handle appd.BtHandle) context.Context {
	u := uuid.New().String()
	ctx := context.WithValue(context.Background(), "APPD", u)
	appd.StoreBT(handle, u)
	return ctx
}

func GetAppdUUID(ctx context.Context) string {
	v := ctx.Value("APPD")
	if v == nil {
		err := errors.New("no APPD Handler found in context")
		log.Println(err)
		return ""
	}

	resp, ok := v.(string)
	if !ok {
		err := fmt.Errorf("objecy %v is not a string pointer", resp)
		log.Println(err)
		return ""
	}

	return resp

}

func GetHandlerFromContext(ctx context.Context) appd.BtHandle {

	return appd.GetBT(GetAppdUUID(ctx))
}

func AddErrorToAppDynamics(ctx context.Context, err error) {
	appd.AddBTError(GetHandlerFromContext(ctx), appd.APPD_LEVEL_ERROR, err.Error(), true)
}
