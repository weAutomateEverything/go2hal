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
	Name      string   `json:"name"`
	PreCheck  *Check   `json:"precheck,omitempty"`
	Actions   []Action `json:"actions"`
	PostCheck *Check   `json:"postcheck,omitempty"`
}

/*
Action object
 */
type Action struct {
	Name string  `json:"name"`
	SearchOption `json:"searchOption"`
	*InputData   `json:"inputdata,omitempty"`
	*ClickButton `json:"clickbutton,omitempty"`
	*ClickLink   `json:"clicklink,omitempty"`
}

/*
SearchOption allows you to specify how you would like to search
 */
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
	Name  string  `json:"name"`
	SearchOption  `json:"searchOption"`
	Value *string `json:"value,omitempty"`
}

/*
CSSSelector Search by CSS
 */
type CSSSelector struct {
}

/*
NameSelector Search by name
 */
type NameSelector struct {
}

/*
TagNameSelector search by tag
 */
type TagNameSelector struct {
}

/*
ClassNameSelector Search by class
 */
type ClassNameSelector struct {
}

/*
IDSelector Search by ID
 */
type IDSelector struct {
}

/*
LinkTextSelector search by link text
 */
type LinkTextSelector struct {
}

/*
PartialLinkTextSelect search by partial link text
 */
type PartialLinkTextSelect struct {
}

/*
XPathSelector search by xpath
 */
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
