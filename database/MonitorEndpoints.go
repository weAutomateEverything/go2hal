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
	Name string
	Endpoint string
	lastChecked time.Time
	lastSuccess time.Time
}

//AddHTMLEndpoint allows for a new endpoint to be added
func AddHTMLEndpoint(name, endpoint string){
	e := HTMLEndpoint{Name:name,Endpoint:endpoint}
	r := htmlEndpoint{HTMLEndpoint: e }
	c := database.C("MonitorHtmlEndpoints")
	c.Insert(r)
}

//GetHTMLEndpoints returns a list of HTML Endpoints
func GetHTMLEndpoints() []HTMLEndpoint{
	c := database.C("MonitorHtmlEndpoints")
	q := c.Find(nil)
	 i,err := q.Count()
	if err != nil {
		log.Println(err)
		return nil
	}
	r := make([]htmlEndpoint,i)
	err = q.All(&r)
	if err != nil {
		log.Println(err)
		return nil
	}

	result := make([]HTMLEndpoint,i)
	for line,x := range r {
		result[line] = x.HTMLEndpoint
	}
	return result
}



