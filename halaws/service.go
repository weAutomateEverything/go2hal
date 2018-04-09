package halaws

import (
	"fmt"
	"github.com/kyokomi/emoji"
	"github.com/weAutomateEverything/aws-sdk-go/aws"
	"github.com/weAutomateEverything/aws-sdk-go/aws/credentials"
	"github.com/weAutomateEverything/aws-sdk-go/aws/session"
	"github.com/weAutomateEverything/aws-sdk-go/service/connect"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"os"
)

type Service interface {
	SendAlert(ctx context.Context, destination string, name string) error
}

type service struct {
	alert alert.Service
}

func NewService(alert alert.Service) Service {
	return &service{alert: alert}
}

func (s *service) SendAlert(ctx context.Context, destination string, name string) error {
	c := credentials.NewEnvCredentials()

	config := aws.Config{Credentials: c, Region: aws.String("us-east-1"), LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)}
	sess, _ := session.NewSession(&config)

	outbound := connect.New(sess, &config)

	req := connect.StartOutboundVoiceContactInput{
		InstanceId:             aws.String(getInstanceID()),
		ContactFlowId:          aws.String(getContactFlowID()),
		DestinationPhoneNumber: aws.String(destination),
		SourcePhoneNumber:      aws.String(getSourcePhoneNumber()),
	}
	output, err := outbound.StartOutboundVoiceContact(&req)
	if err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error invoking alexa to call %v on %v. Error: %v", name, destination, err.Error()))
		return err
	}

	s.alert.SendAlert(ctx, emoji.Sprintf(":phone: HAL has phoned %v on %v. Reference %v ", name, destination, output.ContactId))
	return nil

}

func getInstanceID() string {
	return os.Getenv("AWS_CONNECT_INSTANCE")
}

func getContactFlowID() string {
	return os.Getenv("AWS_CONNECT_FLOW_ID")
}

func getSourcePhoneNumber() string {
	return os.Getenv("AWS_CONNECT_SOURCE_PHONE")
}
