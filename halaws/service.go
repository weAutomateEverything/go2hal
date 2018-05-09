package halaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/kyokomi/emoji"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"os"
	"time"
)

type Service interface {
	SendAlert(ctx context.Context, destination string, name string, variables map[string]string) error
}

type service struct {
	alert alert.Service

	lastcall time.Time
}

func NewService(alert alert.Service) Service {
	return &service{alert: alert}
}

func (s *service) SendAlert(ctx context.Context, destination string, name string, variables map[string]string) error {
	if time.Since(s.lastcall) < time.Duration(30*time.Minute) {
		s.alert.SendAlert(ctx, ":phone: :negative_squared_cross_mark: Not invoking callout since its been less than 30 minutes since the last phone call")
		return nil
	}

	c := credentials.NewEnvCredentials()

	config := aws.Config{Credentials: c, Region: aws.String("us-east-1"), LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody)}
	sess, _ := session.NewSession(&config)

	outbound := connect.New(sess, &config)

	v := map[string]*string{}
	if variables != nil {
		for key, val := range variables {
			v[key] = aws.String(val)
		}
	}

	req := connect.StartOutboundVoiceContactInput{
		InstanceId:             aws.String(getInstanceID()),
		ContactFlowId:          aws.String(getContactFlowID()),
		DestinationPhoneNumber: aws.String(destination),
		SourcePhoneNumber:      aws.String(getSourcePhoneNumber()),
		Attributes:             v,
	}
	output, err := outbound.StartOutboundVoiceContact(&req)
	if err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error invoking alexa to call %v on %v. Error: %v", name, destination, err.Error()))
		return err
	}

	s.lastcall = time.Now()
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
