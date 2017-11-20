package database

import "gopkg.in/mgo.v2/bson"

/*
Selenium Object
 */
type Selenium struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"omitempty"`
	SeleniumServer string
	Name           string
	InitialUrl     string
	Pages          []Page
}

/*
Page object
 */
type Page struct {
	PreCheck  Check
	Actions   []Action
	PostCheck Check
}

/*
Action object
 */
type Action struct {
	Selector    string
	InputData   InputData
	ClickButton ClickButton
	ClickLink   ClickLink
}

/*
InputData Object
 */
type InputData struct {
	Value string
}

/*
ClickButton object
 */
type ClickButton struct {
	Value string
}

/*
ClickLink object
 */
type ClickLink struct {
	Value string
}

/*
Check object
 */
type Check struct {
	Selector string
	Value    string
}

/*
GetAllSeleniumTests returns all selenium tests
 */
func GetAllSeleniumTests() ([]Selenium, error) {
	var result []Selenium
	err := database.C("Selenium").Find(nil).All(&result)
	return result, err
}

/*
AddSelenium adds a test to the database
 */
func AddSelenium(selenium Selenium) error {
	return database.C("Selenium").Insert(selenium)
}
