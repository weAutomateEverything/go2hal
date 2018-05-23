package httpSmoke

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func NewMongoStore(db *mgo.Database) Store {
	return &mongoStore{db}
}

type httpEndpoint struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Endpoint   string
	Method     string
	Parameters []parameters
	Threshold  int
	Chat       uint32

	LastChecked time.Time
	LastSuccess time.Time
	ErrorCount  int
	Passing     bool
	Error       string
}

type parameters struct {
	Name, Value string
}

type Store interface {
	addHTMLEndpoint(endpoint httpEndpoint) error
	getHTMLEndpoints() []httpEndpoint
	getHTMLEndpointsByChat(chat uint32) ([]httpEndpoint, error)
	successfulEndpointTest(endpoint *httpEndpoint) error
	failedEndpointTest(endpoint *httpEndpoint, errorMessage string) error
}

type mongoStore struct {
	mongo *mgo.Database
}

func (s *mongoStore) getHTMLEndpointsByChat(chat uint32) (result []httpEndpoint, err error) {
	c := s.mongo.C("MonitorHtmlEndpoints")
	q := c.Find(bson.M{"chat": chat})
	err = q.All(&result)
	return

}

func (s *mongoStore) addHTMLEndpoint(endpoint httpEndpoint) error {
	c := s.mongo.C("MonitorHtmlEndpoints")
	return c.Insert(endpoint)
}

func (s *mongoStore) getHTMLEndpoints() []httpEndpoint {
	c := s.mongo.C("MonitorHtmlEndpoints")
	q := c.Find(nil)
	i, err := q.Count()
	if err != nil {
		log.Println(err)
		return nil
	}
	r := make([]httpEndpoint, i)
	err = q.All(&r)
	if err != nil {
		log.Println(err)
		return nil
	}
	return r
}

func (s *mongoStore) successfulEndpointTest(endpoint *httpEndpoint) error {
	c := s.mongo.C("MonitorHtmlEndpoints")

	endpoint.LastChecked = time.Now()
	endpoint.LastSuccess = time.Now()
	endpoint.Passing = true
	endpoint.Error = ""
	endpoint.ErrorCount = 0

	err := c.UpdateId(endpoint.ID, endpoint)
	if err != nil {
		return fmt.Errorf("error saving endpoint with success details: %s", err.Error())
	}
	return nil
}

func (s *mongoStore) failedEndpointTest(endpoint *httpEndpoint, errorMessage string) error {
	c := s.mongo.C("MonitorHtmlEndpoints")
	result := httpEndpoint{}
	err := c.FindId(endpoint.ID).One(&result)
	if err != nil {
		return fmt.Errorf("error retreiving endpoint with success details: %s", err.Error())
	}

	result.LastChecked = time.Now()
	result.Passing = false
	result.Error = errorMessage
	result.ErrorCount++

	err = c.UpdateId(endpoint.ID, result)
	if err != nil {
		return fmt.Errorf("error saving endpoint with success details: %s", err.Error())
	}
	return nil
}
