package service

import (
	"fmt"
	"github.com/tebeka/selenium"
	"time"
	"github.com/zamedic/go2hal/database"
	"errors"
)

func init() {
	go func() { runTests() }()
}

var retries = 3

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
				database.IncreaseValue("SELENIUM_RUNS")
				if err != nil {
					SendAlert(fmt.Sprintf("Error executing selenium test for %s. error: %s", test.Name, err.Error()))
				}
			}
		}
		time.Sleep(10 * time.Minute)
	}
}

func doSelenium(item database.Selenium) error {
	if item.SeleniumServer == "" {
		return errors.New("no Selenium Server set")
	}
	if item.Name == "" {
		return errors.New("no script name set")
	}
	if item.InitialURL == "" {
		return errors.New("no initial url set")
	}
	if len(item.Pages) == 0 {
		return errors.New("no pages detected in script")
	}

	var webDriver selenium.WebDriver
	var err error
	caps := selenium.Capabilities(map[string]interface{}{"browserName": "chrome"})
	caps["chrome.switches"] = []string{"--ignore-certificate-errors"}

	if webDriver, err = selenium.NewRemote(caps, item.SeleniumServer); err != nil {
		fmt.Printf("Failed to open session: %s\n", err)
		return err
	}

	defer webDriver.Quit()

	err = webDriver.Get(item.InitialURL)
	if err != nil {
		return handleSeleniumError(item.Name,"Initial Page","Load Page", err, webDriver)
	}

	for _, page := range item.Pages {
		if len(page.Actions) == 0 {
			return fmt.Errorf("no pages found in test %s for page %s", item.Name, page.Name)
		}
		if page.PreCheck != nil {
			if page.PreCheck.SearchPattern == "" {
				return fmt.Errorf("no search pattern found for precheck on test %s, page %s, check %s", item.Name, page.Name, page.PreCheck.Name)
			}
			err = doCheck(page.PreCheck, webDriver)
			if err != nil {
				return handleSeleniumError(item.Name,page.Name,page.PreCheck.Name, err, webDriver)
			}
		}
		for _, action := range page.Actions {
			if action.SearchPattern == "" {
				return fmt.Errorf("no search pattern found for test %s, page %s, action %s", item.Name, page.Name, action.Name)
			}
			elems, err := findElement(action.SearchOption, webDriver)
			if err != nil {
				return handleSeleniumError(item.Name,page.Name,action.Name, err, webDriver)
			}
			elem := elems[0]
			executed := false
			if action.ClickLink != nil {
				elem.Click()
				executed = true
			}
			if action.ClickButton != nil && !executed {
				elem.Click()
				executed = true
			}
			if action.InputData != nil && !executed {
				elem.SendKeys(action.InputData.Value)
				executed = true
			}
			if !executed {
				return fmt.Errorf("no action executed for test %s, page %s, action %s", item.Name, page.Name, action.Name)
			}

		}
		if page.PostCheck != nil {
			if page.PostCheck.SearchPattern == "" {
				return fmt.Errorf("no search pattern found for post check on test %s, page %s, check %s", item.Name, page.Name, page.PostCheck.Name)
			}
			err := doCheck(page.PostCheck, webDriver)
			if err != nil {
				return handleSeleniumError(item.Name,page.Name,page.PostCheck.Name, err, webDriver)
			}
		}
	}
	return nil
}

func doCheck(check *database.Check, driver selenium.WebDriver) error {
	waitfor := func(wb selenium.WebDriver) (bool, error) {

		elems, err := findElement(check.SearchOption, driver)
		if err != nil {
			return false, nil
		}

		for _, elem := range elems {
			dis, err := elem.IsDisplayed()
			if err != nil {
				return false, nil
			}
			if dis {
				if check.Value != nil {
					s, err := elem.Text();
					if err != nil {
						return false, nil
					}
					return *check.Value == s, nil
				}
				return true, nil
			}
		}
		return false, nil
	}
	/*
	Sometimes, selenium throws a EOF error. This is because of some strange behavior within selenium. Lets try 3 times.
	Then Fail.
	 */
	tries:= 0
	for tries < retries {
		err = driver.WaitWithTimeout(waitfor, 10*time.Second)
		if err == nil {
			break;
		}
		tries++
	}
	return err

}

func handleSeleniumError(name, page, action string, err error, driver selenium.WebDriver,) error {
	SendAlert(fmt.Sprintf("%s Selenium Error: %s", name, err.Error()))
	bytes, error := driver.Screenshot()
	if error != nil {
		SendError(error)
		return err
	}
	sendImageToAlertGroup(bytes)
	return err
}

func findElement(action database.SearchOption, driver selenium.WebDriver) ([]selenium.WebElement, error) {
	selector := selenium.ByCSSSelector
	if action.XPathSelector != nil {
		selector = selenium.ByXPATH;
	}
	if action.PartialLinkTextSelect != nil {
		selector = selenium.ByPartialLinkText
	}
	if action.LinkTextSelector != nil {
		selector = selenium.ByLinkText
	}
	if action.IDSelector != nil {
		selector = selenium.ByID
	}
	if action.ClassNameSelector != nil {
		selector = selenium.ByCSSSelector
	}
	if action.NameSelector != nil {
		selector = selenium.ByName
	}
	if action.TagNameSelector != nil {
		selector = selenium.ByTagName
	}
	if action.Multiple {
		return driver.FindElements(selector, action.SearchPattern)
	}
	elem, err := driver.FindElement(selector, action.SearchPattern)
	return []selenium.WebElement{elem}, err

}
