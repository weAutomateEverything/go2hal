package service

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
	"github.com/zamedic/go2hal/database"
)

func init(){
	go func() {runTests()}()
}

/*
TestSelenium tests a selenium endpoint and adds it to the database.
 */
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

func runTests() {
	for true {
		tests, err := database.GetAllSeleniumTests()
		if err != nil {
			SendError(err)
		} else {
			for _, test := range tests {
				err = doSelenium(test)
				if err != nil {
					SendAlert(fmt.Sprintf("Error executing selenium test for %s. error: %s", test.Name, err.Error()))
				}
			}
		}
		time.Sleep(10 * time.Minute)
	}
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
	return driver.WaitWithTimeout(waitfor, 10*time.Second)
}
