package http

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
	"fmt"
	"gopkg.in/mgo.v2"
)

type Store interface {
	addHTMLEndpoint(endpoint httpEndpoint)

	getHTMLEndpoints() []httpEndpoint

	successfulEndpointTest(endpoint *httpEndpoint) error

	failedEndpointTest(endpoint *httpEndpoint, errorMessage string) error
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(db *mgo.Database)Store{
	return &mongoStore{db}
}

type httpEndpoint struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Endpoint   string
	Method     string
	Parameters []parameters
	Threshold  int

	LastChecked time.Time
	LastSuccess time.Time
	ErrorCount  int
	Passing     bool
	Error       string
}

type parameters struct {
	Name, Value string
}

func (s *mongoStore)addHTMLEndpoint(endpoint httpEndpoint) {
	c := s.mongo.C("MonitorHtmlEndpoints")
	c.Insert(endpoint)
}

func (s *mongoStore)getHTMLEndpoints() []httpEndpoint {
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


func (s *mongoStore)successfulEndpointTest(endpoint *httpEndpoint) error {
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


func (s *mongoStore)failedEndpointTest(endpoint *httpEndpoint, errorMessage string) error {
	c := s.mongo.C("MonitorHtmlEndpoints")
	result := httpEndpoint{}
	err := c.FindId(endpoint.ID).One(&result);
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
