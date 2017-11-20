package service

import (
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"time"
	"github.com/zamedic/go2hal/database"
)

func TestSelenium2() {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	if webDriver, err = selenium.NewRemote(caps, "http://card-selenium-chrome-dev.chop.standardbank.co.za/wd/hub"); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return
	}
	defer webDriver.Quit()

	err = webDriver.Get("https://dinerspbweb-dev.standardbank.co.za/")
	elem, err := webDriver.FindElement(selenium.ByName, "username")
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = elem.SendKeys("c1592023")
	if err != nil {
		log.Println(err.Error())
		return
	}
	elem, err = webDriver.FindElement(selenium.ByName, "password")
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = elem.SendKeys("trendweb")
	if err != nil {
		log.Println(err.Error())
		return
	}

	elem, err = webDriver.FindElement(selenium.ByCSSSelector, ".md-primary")
	if err != nil {
		log.Println(err.Error())
		return
	}

	elem.Click()
	loginSuccess := func(wb selenium.WebDriver) (bool, error) {
		elem, err := wb.FindElement(selenium.ByName, "contactPerson")
		if err != nil {
			return false, nil
		}
		return elem.IsDisplayed()
	}

	err = webDriver.WaitWithTimeout(loginSuccess, 10*time.Second)
	if err != nil {
		log.Println(err.Error())
		return
	}

	elem, err = webDriver.FindElement(selenium.ByTagName, "h2")
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(elem.Text())

}

func TestSelenium(item database.Selenium) error {
	err := doSelenium(item)
	if err != nil {
		return err
	}

	err = database.AddSelenium(item)
	if err != nil {
		return err
	}
	return nil
}

func doSelenium(item database.Selenium) error {
	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	if webDriver, err = selenium.NewRemote(caps, item.SeleniumServer); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return err
	}
	defer webDriver.Quit()

	err = webDriver.Get(item.InitialUrl)
	if err != nil {
		fmt.Printf("Failed to open initial page %s\n", err)
		return err
	}

	for _, page := range item.Pages {
		if page.PreCheck.Selector != "" {
			err = doCheck(page.PreCheck, webDriver)
			if err != nil {
				return err
			}
		}
		for _, action := range page.Actions {
			elem, err := webDriver.FindElement(selenium.ByCSSSelector, action.Selector)
			if err != nil {
				fmt.Printf("Failed to find element: %s\n", err)
				return err
			}
			if action.ClickLink.Value != "" {
				elem.Click()
			}
			if action.ClickButton.Value != "" {
				elem.Click()
			}
			if action.InputData.Value != "" {
				elem.SendKeys(action.InputData.Value)
			}

		}

		if page.PostCheck.Selector != "" {
			err := doCheck(page.PostCheck, webDriver)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func doCheck(check database.Check, driver selenium.WebDriver) error {
	waitfor := func(wb selenium.WebDriver) (bool, error) {
		elems, err := wb.FindElements(selenium.ByCSSSelector, check.Selector)
		if err != nil {
			return false, nil
		}

		for _, elem := range elems {
			dis, err := elem.IsDisplayed()
			if err != nil {
				return false, nil
			}
			if (dis) {
				if check.Value != "" {
					s, err := elem.Text();
					if err != nil {
						return false, nil
					}
					return check.Value == s, nil
				}
				return true, nil
			}
		}
		return false, nil
	}

	return driver.WaitWithTimeout(waitfor,10*time.Second)

}
