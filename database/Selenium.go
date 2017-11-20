package database

import "gopkg.in/mgo.v2/bson"

type Selenium struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"omitempty"`
	SeleniumServer string
	Name           string
	InitialUrl     string
	Pages          []Page
}

type Page struct {
	PreCheck  Check
	Actions   []Action
	PostCheck Check
}

type Action struct {
	Selector    string
	InputData   InputData
	ClickButton ClickButton
	ClickLink   ClickLink
}

type InputData struct {
	Value string
}

type ClickButton struct {
	Value string
}

type ClickLink struct {
	Value string
}

type Check struct {
	Selector string
	Value    string
}

func GetAllSeleniumTests() ([]Selenium, error) {
	var result []Selenium
	err := database.C("Selenium").Find(nil).All(&result)
	return result, err
}

func AddSelenium(selenium Selenium) error {
	return database.C("Selenium").Insert(selenium)
}
