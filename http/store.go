package http

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
	"fmt"
	"gopkg.in/mgo.v2"
)

type Store interface {
	//AddHTMLEndpoint allows for a new endpoint to be added
	AddHTMLEndpoint(endpoint HTTPEndpoint)

	//GetHTMLEndpoints returns a list of HTML Endpoints
	GetHTMLEndpoints() []HTTPEndpoint

	/*
	SuccessfulEndpointTest will update the mongo element with the ID with the latest details to show it passed successfully
	*/
	SuccessfulEndpointTest(endpoint *HTTPEndpoint) error

	/*
	FailedEndpointTest will update the mongo element with the failed details
	*/
	FailedEndpointTest(endpoint *HTTPEndpoint, errorMessage string) error
}

type mongoStore struct {
	mongo mgo.Database
}

type HTTPEndpoint struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Endpoint   string
	Method     string
	Parameters []Parameters
	Threshold  int

	LastChecked time.Time
	LastSuccess time.Time
	ErrorCount  int
	Passing     bool
	Error       string
}

type Parameters struct {
	Name, Value string
}

func (s *mongoStore)AddHTMLEndpoint(endpoint HTTPEndpoint) {
	c := s.mongo.C("MonitorHtmlEndpoints")
	c.Insert(endpoint)
}

func (s *mongoStore)GetHTMLEndpoints() []HTTPEndpoint {
	c := s.mongo.C("MonitorHtmlEndpoints")
	q := c.Find(nil)
	i, err := q.Count()
	if err != nil {
		log.Println(err)
		return nil
	}
	r := make([]HTTPEndpoint, i)
	err = q.All(&r)
	if err != nil {
		log.Println(err)
		return nil
	}
	return r
}


func (s *mongoStore)SuccessfulEndpointTest(endpoint *HTTPEndpoint) error {
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


func (s *mongoStore)FailedEndpointTest(endpoint *HTTPEndpoint, errorMessage string) error {
	c := s.mongo.C("MonitorHtmlEndpoints")
	result := HTTPEndpoint{}
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
