package config

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Store interface {
	/*
	SaveCalloutDetails saves the new callout configs to the db
	*/
	SaveCalloutDetails(URL string) error

	/*
	GetCalloutDetails returns the callout details
	*/
	GetCalloutDetails() (*CallOut, error)

	/*
	SaveJiraDetails saves the JIRA details
	*/
	SaveJiraDetails(url, template, defaultUser string) error

	/*
	GetJiraDetails returns the current JIRA details
	*/
	GetJiraDetails() (*Jira, error)
}

type mongoStore struct {
	mongo mgo.Database
}

type config struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	*CallOut
	*Jira
}

/*
CallOut details
 */
type CallOut struct {
	URL string
}

/*
Jira details
 */
type Jira struct {
	URL         string
	Template    string
	DefaultUser string
}


func (s *mongoStore)SaveCalloutDetails(URL string) error {
	c, err := getConfig(s.mongo)
	if err != nil {
		return err
	}
	c.CallOut = &CallOut{URL: URL}
	return saveConfig(s.mongo,c)
}


func (s *mongoStore)GetCalloutDetails() (*CallOut, error) {
	c, err := getConfig(s.mongo)
	if err != nil {
		return nil, err
	}
	return c.CallOut, nil
}


func (s *mongoStore)SaveJiraDetails(url, template, defaultUser string) error {
	c, err := getConfig(s.mongo)
	if err != nil {
		return err
	}
	c.Jira = &Jira{URL: url, Template: template, DefaultUser: defaultUser}
	return saveConfig(s.mongo,c)
}


func (s *mongoStore)GetJiraDetails() (*Jira, error) {
	c, err := getConfig(s.mongo)
	if err != nil {
		return nil, err
	}
	return c.Jira, nil
}

func getConfig(database mgo.Database) (*config, error) {
	c := database.C("Config")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return &config{}, nil
	}
	cfg := config{}
	err = q.One(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil

}

func saveConfig(database mgo.Database,config *config) error {
	c := database.C("Config")
	if config.ID == "" {
		return c.Insert(config)
	}
	return c.UpdateId(config.ID, config)
}
