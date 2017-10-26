package database

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type appDynamics struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	MqEndpoints []MqEndpoint
}

/*
MqEndpoint Object
 */
type MqEndpoint struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Endpoint   string
	MetricPath string
}

/*
AddMqEndpoint will add an MQ endpoint to be monitored
 */
func AddMqEndpoint(name, endpoint string, metricPath string) {
	var mq = MqEndpoint{Endpoint: endpoint, MetricPath: metricPath, Name: name}
	appd := getAppDynamics()
	appd.MqEndpoints = append(appd.MqEndpoints, mq)
	c := database.C("appDynamics")
	err := c.UpdateId(appd.ID, appd)
	if err != nil {
		log.Printf("Error saving to db %s", err)
	}
}

/*
GetMQEndponts will return a list of MQ endpoints configured
 */
func GetMQEndponts() []MqEndpoint {
	appd := getAppDynamics()
	return appd.MqEndpoints
}

func getAppDynamics() appDynamics {
	c := database.C("appDynamics")
	i, err := c.Count()
	if err != nil || i == 0 {
		appd := appDynamics{}
		c.Insert(appd)
		return appd
	}
	var appd []appDynamics
	c.Find(nil).All(&appd)
	return appd[0]
}
