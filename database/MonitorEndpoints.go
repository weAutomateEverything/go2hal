package database

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

type htmlEndpoint struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	HtmlEndpoint
}

type HtmlEndpoint struct {
	Name string
	Endpoint string
	lastChecked time.Time
	lastSuccess time.Time
}

func AddHtmlEndpoint(name, endpoint string){
	e := HtmlEndpoint{Name:name,Endpoint:endpoint}
	r := htmlEndpoint{HtmlEndpoint: e }
	c := database.C("MonitorHtmlEndpoints")
	c.Insert(r)
}

func GetHtmlEndpoints() []HtmlEndpoint{
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

	result := make([]HtmlEndpoint,i)
	for line,x := range r {
		result[line] = x.HtmlEndpoint
	}

	return result

}



