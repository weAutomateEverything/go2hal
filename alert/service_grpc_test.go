package alert

import (
	"github.com/golang/mock/gomock"
	"github.com/weAutomateEverything/go2hal/telegram"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"testing"
)

func TestAlertGrpc(t *testing.T) {
	ctrl := gomock.NewController(t)

	telegramService := telegram.NewMockService(ctrl)
	telegramStore := telegram.NewMockStore(ctrl)

	alertService := NewService(telegramService, telegramStore)

	mockServer := NewMockAlert_RequestReplyServer(ctrl)

	a := AlertRequest{
		Message:       "Hello World",
		CorrelationId: "12345",
		GroupID:       54321,
	}

	mockServer.EXPECT().Recv().Times(1).Return(&a, nil)
	mockServer.EXPECT().Context().Times(1).Return(context.Background())

	telegramStore.EXPECT().GetRoomKey(uint32(54321)).Times(1).Return(int64(11223344), nil)

	telegramService.EXPECT().SendMessageWithCorrelation(context.Background(), int64(11223344), "Hello World", 0, "12345").Return(1111, nil)

	grpcService := NewGrpcService(alertService, telegramStore)

	g := grpc.NewServer()
	RegisterAlertServer(g, grpcService)
	reflection.Register(g)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		g.Serve(ln)
	}()

	con, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}

	c := NewAlertClient(con)

	client, err := c.RequestReply(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	err = client.Send(&a)
	if err != nil {
		t.Fatal(err)
	}

}
