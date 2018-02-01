package halSelenium

import (
	"github.com/tebeka/selenium"
	"github.com/zamedic/go2hal/alert"
	"github.com/pkg/errors"
	"fmt"
	selenium2 "github.com/zamedic/go2hal/seleniumTests"
)

type Service interface {
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

func NewChromeService(service alert.Service, server string) Service {

	s := &chromeService{alert: service}
	err := s.newClient(server)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *chromeService) Driver() selenium.WebDriver {
	return s.driver
}

func (s *chromeService) newClient(seleniumServer string) error {
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
			err := s.alert.SendImageToHeartbeatGroup(image)
			if err != nil {
				return err
			}
		} else {
			err := s.alert.SendImageToAlertGroup(image)
			if err != nil {
				return err
			}
		}
	}

	if internalError {
		s.alert.SendError(errors.New(message))
	} else {
		err := s.alert.SendAlert(message)
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
