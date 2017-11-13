package database

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

type htmlEndpoint struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	HTMLEndpoint
}

//HTMLEndpoint reprisents a HTTP Endpoint that the system will monitor
type HTMLEndpoint struct {
	IdString    string `bson:"-"`
	Name        string
	Endpoint    string
	lastChecked time.Time
	lastSuccess time.Time
	Passing     bool
	Error       string
}

//AddHTMLEndpoint allows for a new endpoint to be added
func AddHTMLEndpoint(name, endpoint string) {
	e := HTMLEndpoint{Name: name, Endpoint: endpoint}
	r := htmlEndpoint{HTMLEndpoint: e}
	c := database.C("MonitorHtmlEndpoints")
	c.Insert(r)
}

//GetHTMLEndpoints returns a list of HTML Endpoints
func GetHTMLEndpoints() []HTMLEndpoint {
	c := database.C("MonitorHtmlEndpoints")
	q := c.Find(nil)
	i, err := q.Count()
	if err != nil {
		log.Println(err)
		return nil
	}
	r := make([]htmlEndpoint, i)
	err = q.All(&r)
	if err != nil {
		log.Println(err)
		return nil
	}

	result := make([]HTMLEndpoint, i)
	for line, x := range r {
		result[line] = x.HTMLEndpoint
		result[line].IdString = x.ID.String()
	}
	return result
}

/*
SuccessfulEndpointTest will update the mongo element with the ID with the latest details to show it passed successfully
 */
func SuccessfulEndpointTest(id string) {
	c := database.C("MonitorHtmlEndpoints")
	result := htmlEndpoint{}
	err := c.FindId(id).One(&result);
	if err != nil {
		log.Printf("Error retreiving endpoint with success details: %s", err.Error())
		return
	}

	result.lastChecked = time.Now()
	result.lastSuccess = time.Now()
	result.Passing = true
	result.Error = ""

	err = c.UpdateId(id, result)
	if err != nil {
		log.Printf("Error saving endpoint with success details: %s", err.Error())
	}
}

/*
FailedEndpointTest will update the mongo element with the failed details
 */
func FailedEndpointTest(id,errorMessage string ) {
	c := database.C("MonitorHtmlEndpoints")
	result := htmlEndpoint{}
	err := c.FindId(id).One(&result);
	if err != nil {
		log.Printf("Error retreiving endpoint with success details: %s", err.Error())
		return
	}

	result.lastChecked = time.Now()
	result.Passing = false
	result.Error = errorMessage

	err = c.UpdateId(id, result)
	if err != nil {
		log.Printf("Error saving endpoint with success details: %s", err.Error())
	}

}
