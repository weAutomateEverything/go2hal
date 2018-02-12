package remoteTelegramCommands

import (
	"google.golang.org/grpc"
	"os"
)

func NewRemoteCommandClientService() RemoteCommandClient {

	conn, err := grpc.Dial(getURL(""), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return NewRemoteCommandClient(conn)
}

func getURL(namespace string) string {
	u := getAlertUrl()
	if u != "" {
		return u
	}
	u = "hal"
	if namespace != "" {
		u = u + "." + namespace
	}
	u = u + ":8080"
	return u
}

func getAlertUrl() string {
	return os.Getenv("HAL_GRPC_ENDPOINT")
}
