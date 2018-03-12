package halaws

import (
	"github.com/weAutomateEverything/aws-sdk-go/aws"
	"github.com/weAutomateEverything/aws-sdk-go/aws/credentials"
	"github.com/weAutomateEverything/aws-sdk-go/aws/session"
	"github.com/weAutomateEverything/aws-sdk-go/service/connect"
	"os"
)

type Service interface {
	SendAlert(destination string)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) SendAlert(destination string) {
	c := credentials.NewEnvCredentials()

	config := aws.Config{Credentials: c, Region: aws.String("us-east-1"), LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)}
	sess, _ := session.NewSession(&config)

	outbound := connect.New(sess, &config)

	req := connect.StartOutboundVoiceContactInput{
		InstanceId:             aws.String(getInstanceID()),
		ContactFlowId:          aws.String(getContactFlowID()),
		SourcePhoneNumber:      aws.String(getSourcePhoneNumber()),
		DestinationPhoneNumber: aws.String(destination),
		Attributes:             &map[string]string{},
	}
	outbound.StartOutboundVoiceContact(&req)
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
