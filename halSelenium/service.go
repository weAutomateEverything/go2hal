package halSelenium

import (
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/pkg/errors"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	selenium2 "github.com/zamedic/go2hal/seleniumTests"
	"golang.org/x/net/context"
)

type Service interface {
	NewClient(seleniumServer string) error
	HandleSeleniumError(internal bool, err error)
	Driver() selenium.WebDriver

	ClickByClassName(cn string)
	ClickByXPath(xp string)
	ClickByCSSSelector(cs string)

	WaitFor(findBy, selector string)
}

type chromeService struct {
	alert  alert.Service
	driver selenium.WebDriver
}

func NewChromeService(service alert.Service) Service {
	fieldKeys := []string{"method"}
	s := newChromeService(service)
	s = NewInstrumentService(kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "halSelenium",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys),
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "halSelenium",
			Name:      "error_count",
			Help:      "Number of errors encountered.",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "halSelenium",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys), s)
	return s
}

func newChromeService(service alert.Service) Service {

	s := &chromeService{alert: service}
	return s
}

func (s *chromeService) Driver() selenium.WebDriver {
	return s.driver
}

func (s *chromeService) NewClient(seleniumServer string) error {
	driver, err := selenium2.NewChromeClient(seleniumServer)
	if err != nil {
		return err
	}
	s.driver = driver
	return nil
}

func (s chromeService) HandleSeleniumError(internal bool, err error) {
	msg := err.Error()
	if s.driver == nil {
		s.sendError(msg, nil, internal)
		return
	}
	bytes, error := s.driver.Screenshot()
	if error != nil {
		// Couldnt get a screenshot - lets end the original error
		s.sendError(msg, nil, internal)
		return
	}
	s.sendError(msg, bytes, internal)
}

func (s chromeService) ClickByClassName(cn string) {
	item, err := s.driver.FindElement(selenium.ByClassName, cn)
	if err != nil {
		panic(err)
	}

	err = item.Click()
	if err != nil {
		panic(err)
	}
}

func (s chromeService) ClickByXPath(xp string) {
	item, err := s.driver.FindElement(selenium.ByXPATH, xp)
	if err != nil {
		panic(err)
	}

	err = item.Click()
	if err != nil {
		panic(err)
	}
}

func (s chromeService) ClickByCSSSelector(cs string) {
	item, err := s.driver.FindElement(selenium.ByCSSSelector, cs)
	if err != nil {
		panic(err)
	}

	err = item.Click()
	if err != nil {
		panic(err)
	}
}

func (s *chromeService) WaitFor(findBy, selector string) {

	e := s.driver.Wait(func(wb selenium.WebDriver) (bool, error) {

		elem, err := wb.FindElement(findBy, selector)
		if err != nil {
			return false, nil
		}
		r, err := elem.IsDisplayed()
		return r, nil
	})
	if e != nil {
		panic(e)
	}
}

func (s chromeService) sendError(message string, image []byte, internalError bool) error {

	if image != nil {
		if internalError {
			err := s.alert.SendImageToHeartbeatGroup(context.TODO(), image)
			if err != nil {
				return err
			}
		} else {
			err := s.alert.SendImageToAlertGroup(context.TODO(), image)
			if err != nil {
				return err
			}
		}
	}

	if internalError {
		s.alert.SendError(context.TODO(), errors.New(message))
	} else {
		err := s.alert.SendAlert(context.TODO(), message)
		if err != nil {
			return err
		}
	}
	return nil

}

type SeleniumnError struct {
	Internal bool
	Message  error
}

func (e *SeleniumnError) Error() string {
	return fmt.Sprintf("internal: %v, message: %v", e.Internal, e.Message)
}
