package database

import (
	"gopkg.in/mgo.v2/bson"
)

type config struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	*CallOut
	*Jira
}

type CallOut struct {
	URL string
}

type Jira struct {
	URL         string
	Template    string
	DefaultUser string
}

/*
SaveCalloutDetails saves the new callout configs to the db
 */
func SaveCalloutDetails(URL string) error {
	c, err := getConfig()
	if err != nil {
		return err
	}
	c.CallOut = &CallOut{URL: URL}
	return saveConfig(c)
}

/*
GetCalloutDetails returns the callout details
 */
func GetCalloutDetails() (*CallOut, error) {
	c, err := getConfig()
	if err != nil {
		return nil, err
	}
	return c.CallOut, nil
}

/*
SaveJiraDetails saves the JIRA details
 */
func SaveJiraDetails(url, template, defaultUser string) error {
	c , err:= getConfig()
	if err != nil {
		return err
	}
	c.Jira = &Jira{URL: url, Template: template,DefaultUser:defaultUser}
	return saveConfig(c)
}

/*
GetJiraDetails returns the current JIRA details
 */
func GetJiraDetails() (*Jira, error) {
	c, err := getConfig()
	if err != nil {
		return nil, err
	}
	return c.Jira, err
}

func getConfig() (*config, error){
	c := database.C("Config")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return &config{}, nil
	} else {
		c := config{}
		err := q.One(&c)
		if err != nil {
			return nil,err
		}
		return &c, nil
	}
}

func saveConfig(config *config) error {
	c := database.C("Config")
	if config.ID == "" {
		return  c.Insert(config)
	} else {
		return c.UpdateId(config.ID, config)
	}

}
