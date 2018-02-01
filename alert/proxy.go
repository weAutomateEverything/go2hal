package alert

import (
	"context"
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

	sendAlertEndpoint                 endpoint.Endpoint
	sendNonTechnicalAlertEndpoint     endpoint.Endpoint
	sendHeartbeatGroupAlertEndpoint   endpoint.Endpoint
	sendImageToAlertGroupEndpoint     endpoint.Endpoint
	sendImageToHeartbeatGroupEndpoint endpoint.Endpoint
	sendErrorEndpoint                 endpoint.Endpoint
}

/*
NewKubernetesAlertProxy will return an alert service that is actually a HTTP Proxy into the kubertes service
*/
func NewKubernetesAlertProxy(namespace string) Service {
	alert := makeAlertKubernetesHTTPProxy(namespace)
	alert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alert)
	alert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(alert)

	alertImage := makeAlertKubernetesSendImageToAlertGroupHTTPProxy(namespace)
	alertImage = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertImage)
	alertImage = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(alertImage)

	nonTechAlert := makeAlertKubernetesSendNonTechnicalAlertHTTPProxy(namespace)
	nonTechAlert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(nonTechAlert)
	nonTechAlert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(nonTechAlert)

	heartbeatAlert := makeAlertKubernetesSendHeartbeatGroupAlertHTTPProxy(namespace)
	heartbeatAlert = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(heartbeatAlert)
	heartbeatAlert = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(heartbeatAlert)

	heartbeatImage := makeAlertKubernetesSendImageToHeartbeatGroupHTTPProxy(namespace)
	heartbeatImage = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(heartbeatImage)
	heartbeatImage = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(heartbeatImage)

	alertError := makeAlertKubernetesSendErrorHTTPProxy(namespace)
	alertError = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(alertError)
	alertError = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))(alertError)

	return &alertKubernetesProxy{ctx: context.Background(), sendAlertEndpoint: alert, sendErrorEndpoint: alertError,
		sendHeartbeatGroupAlertEndpoint: heartbeatAlert, sendImageToAlertGroupEndpoint: alertImage,
		sendImageToHeartbeatGroupEndpoint: heartbeatImage, sendNonTechnicalAlertEndpoint: nonTechAlert}

}
func (s *alertKubernetesProxy) SendAlert(message string) error {
	_, err := s.sendAlertEndpoint(s.ctx, message)
	return err
}

func (s *alertKubernetesProxy) SendNonTechnicalAlert(message string) error {
	_, err := s.sendNonTechnicalAlertEndpoint(s.ctx, message)
	return err
}
func (s *alertKubernetesProxy) SendHeartbeatGroupAlert(message string) error {
	_, err := s.sendHeartbeatGroupAlertEndpoint(s.ctx, message)
	return err
}
func (s *alertKubernetesProxy) SendImageToAlertGroup(image []byte) error {
	_, err := s.sendImageToAlertGroupEndpoint(s.ctx, image)
	return err
}
func (s *alertKubernetesProxy) SendImageToHeartbeatGroup(image []byte) error {
	_, err := s.sendImageToHeartbeatGroupEndpoint(s.ctx, image)
	return err
}
func (s *alertKubernetesProxy) SendError(err error) {
	s.sendErrorEndpoint(s.ctx, err)
}

func makeAlertKubernetesHTTPProxy(namespace string) endpoint.Endpoint {

	return http.NewClient(
		"POST",
		getURL(namespace, "alert"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
	).Endpoint()

}

func makeAlertKubernetesSendNonTechnicalAlertHTTPProxy(namespace string) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		getURL(namespace, "alert/business"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
	).Endpoint()
}

func makeAlertKubernetesSendHeartbeatGroupAlertHTTPProxy(namespace string) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		getURL(namespace, "alert/heartbeat"),
		gokit.EncodeRequest,
		gokit.DecodeResponse,
	).Endpoint()
}

func makeAlertKubernetesSendImageToAlertGroupHTTPProxy(namespace string) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		getURL(namespace, "alert/image"),
		gokit.EncodeToBase64,
		gokit.DecodeResponse,
	).Endpoint()
}

func makeAlertKubernetesSendImageToHeartbeatGroupHTTPProxy(namespace string) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		getURL(namespace, "alert/heartbeat/image"),
		gokit.EncodeToBase64,
		gokit.DecodeResponse,
	).Endpoint()
}

func makeAlertKubernetesSendErrorHTTPProxy(namespace string) endpoint.Endpoint {
	return http.NewClient(
		"POST",
		getURL(namespace, "alert/image"),
		gokit.EncodeErrorRequest,
		gokit.DecodeResponse,
	).Endpoint()
}

func getURL(namespace, uri string) (*url.URL) {
	u := "http://hal"
	if namespace != "" {
		u = u + "." + namespace
	}
	u = u + "/" + uri

	ur, err := url.Parse(u)

	if err != nil {
		panic(err)
	}

	return ur
}
