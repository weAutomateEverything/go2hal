package alert

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport/http"
	"github.com/sony/gobreaker"
	"github.com/zamedic/go2hal/gokit"
	"golang.org/x/time/rate"
	"net/url"
	"time"
)

type alertKubernetesProxy struct {
	ctx context.Context

	sendAlertEndpoint endpoint.Endpoint
}

/*
NewKubernetesAlertProxy will return an alert service that is actually a HTTP Proxy into the kubertes service
 */
func NewKubernetesAlertProxy(namespace string) Service {
	e := makeAlertKubernetesHTTPProxy(namespace)
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
	e = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(e)

	return &alertKubernetesProxy{ctx: context.Background(), sendAlertEndpoint: e}

}
func (s *alertKubernetesProxy) SendAlert(message string) error {
	_, err := s.sendAlertEndpoint(s.ctx, message)
	return err
}

func (s *alertKubernetesProxy) SendNonTechnicalAlert(message string) error {
	return nil
}
func (s *alertKubernetesProxy) SendHeartbeatGroupAlert(message string) error {
	return nil
}
func (s *alertKubernetesProxy) SendImageToAlertGroup(image []byte) error {
	return nil
}
func (s *alertKubernetesProxy) SendImageToHeartbeatGroup(image []byte) error {
	return nil
}
func (s *alertKubernetesProxy) SendError(err error) {

}

func makeAlertKubernetesHTTPProxy(namespance string) endpoint.Endpoint {
	u := ""
	if namespance == "" {
		u = "http://hal/alert"
	} else {
		u = fmt.Sprintf("http://hal.%v/alert", namespance)
	}

	ur, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return http.NewClient(
		"POST",
		ur,
		gokit.EncodeRequest,
		gokit.DecodeResponse,
	).Endpoint()

}
