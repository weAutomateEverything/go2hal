package halaws

import (
	"crypto/tls"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/kyokomi/emoji"
	"github.com/weAutomateEverything/go2hal/alert"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"time"
)

type Service interface {
	SendAlert(ctx context.Context, chatId uint32, destination string, name string, variables map[string]string) error
}

type service struct {
	alert alert.Service

	lastcall map[uint32]time.Time
}

func NewService(alert alert.Service) Service {
	s := &service{alert: alert}
	s.lastcall = make(map[uint32]time.Time)
	return s
}

func (s *service) SendAlert(ctx context.Context, chatId uint32, destination string, name string, variables map[string]string) error {
	if !s.checkCallout(ctx, chatId) {
		return nil
	}

	c := credentials.NewEnvCredentials()

	client := http.DefaultClient
	transport := http.DefaultTransport
	transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client.Transport = transport
	config := aws.Config{Credentials: c, Region: aws.String(os.Getenv("AWS_REGION")), LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody), HTTPClient: client}
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
	output, err := outbound.StartOutboundVoiceContactWithContext(ctx, &req)
	if err != nil {
		s.alert.SendError(ctx, fmt.Errorf("error invoking alexa to call %v on %v. Error: %v", name, destination, err.Error()))
		return err
	}

	s.lastcall[chatId] = time.Now()
	s.alert.SendAlert(ctx, chatId, emoji.Sprintf(":phone: HAL has phoned %v on %v. Reference %v ", name, destination, output.ContactId))
	return nil

}

func (s service) checkCallout(ctx context.Context, chatid uint32) bool {
	t, ok := s.lastcall[chatid]
	if ok {
		if time.Since(t) < time.Duration(30*time.Minute) {
			s.alert.SendAlert(ctx, chatid, emoji.Sprintf(":phone: :negative_squared_cross_mark: Not invoking callout since its been less than 30 minutes since the last phone call"))
			return false
		}
	}
	return true
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
