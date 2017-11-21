package database

import "gopkg.in/mgo.v2/bson"

/*
Selenium Object
 */
type Selenium struct {
	ID             bson.ObjectId `bson:"_id,omitempty" json:"omitempty"`
	SeleniumServer string
	Name           string
	InitialURL     string
	Pages          []Page
}

/*
Page object
 */
type Page struct {
	PreCheck  *Check   `json:"precheck,omitempty"`
	Actions   []Action `json:"actions"`
	PostCheck *Check   `json:"postcheck,omitempty"`
}

/*
Action object
 */
type Action struct {
	Selector    string       `json:"selector"`
	InputData   *InputData   `json:"inputdata,omitempty"`
	ClickButton *ClickButton `json:"clickbutton,omitempty"`
	ClickLink   *ClickLink   `json:"clicklink,omitempty"`
}

/*
InputData Object
 */
type InputData struct {
	Value string `json:"value"`
}

/*
ClickButton object
 */
type ClickButton struct {
}

/*
ClickLink object
 */
type ClickLink struct {
}

/*
Check object
 */
type Check struct {
	Selector string  `json:"selector"`
	Value    *string `json:"value,omitempty"`
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
