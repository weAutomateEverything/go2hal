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
	SearchOption `json:"searchOption"`
	*InputData   `json:"inputdata,omitempty"`
	*ClickButton `json:"clickbutton,omitempty"`
	*ClickLink   `json:"clicklink,omitempty"`
}

type SearchOption struct {
	Multiple      bool     `json:"multiple"`
	SearchPattern string   `json:"searchPattern"`
	*CSSSelector           `json:"CSSSelector,omitempty"`
	*NameSelector          `json:"nameSelector,omitempty"`
	*TagNameSelector       `json:"tagNameSelector,omitempty"`
	*ClassNameSelector     `json:"classNameSelector,omitempty"`
	*IDSelector            `json:"IDSelector,omitempty"`
	*LinkTextSelector      `json:"linkTextSelector,omitempty"`
	*PartialLinkTextSelect `json:"partialLinkTextSelect,omitempty"`
	*XPathSelector         `json:"XPathSelector,omitempty"`
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
	SearchOption  `json:"searchOption"`
	Value *string `json:"value,omitempty"`
}

type CSSSelector struct {
}

type NameSelector struct {
}

type TagNameSelector struct {
}

type ClassNameSelector struct {
}

type IDSelector struct {
}

type LinkTextSelector struct {
}

type PartialLinkTextSelect struct {
}

type XPathSelector struct {
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
