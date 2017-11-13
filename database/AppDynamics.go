package database

import (
	"gopkg.in/mgo.v2/bson"
	"errors"
)

type appDynamics struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Endpoint    string
	MqEndpoints []MqEndpoint
}

/*
MqEndpoint Object
 */
type MqEndpoint struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Name       string
	Application   string
	MetricPath string
}

/*
AddAppDynamicsEndpoint will add a app dynamics endpoint to the mongo DB if it doesnt exist. If it exists,. it will
update it.
*/
func AddAppDynamicsEndpoint(endpoint string) error{
	a := appDynamics{Endpoint: endpoint}
	b, err := GetAppDynamics()
	c := database.C("appDynamics")

	if err == nil {
		a.ID = b.ID
		a.MqEndpoints = b.MqEndpoints
		err = c.UpdateId(a.ID,a)
		if err != nil {
			return err
		}
	} else {
		c.Insert(a)
	}
	return nil
}


/*
AddMqEndpoint will add an MQ endpoint to be monitored
 */
func AddMqEndpoint(name, application string, metricPath string) error {
	var mq = MqEndpoint{Application: application, MetricPath: metricPath, Name: name}
	appd,err := GetAppDynamics()
	if err != nil {
		return err
	}

	appd.MqEndpoints = append(appd.MqEndpoints, mq)
	c := database.C("appDynamics")
	err = c.UpdateId(appd.ID, appd)
	if err != nil {
		return err
	}
	return nil
}


/*
GetAppDynamics wll return the app dynamics object in the ob, Else, error if nothing exists.
 */
func GetAppDynamics() (appDynamics, error) {
	c := database.C("appDynamics")
	i, err := c.Count()
	if err != nil || i == 0 {
		return appDynamics{}, errors.New("no app dynamics config set")
	}
	a := appDynamics{}
	err = c.Find(nil).One(&a)
	return a,err
}
