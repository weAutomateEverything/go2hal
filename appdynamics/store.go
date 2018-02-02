package appdynamics

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	/*
		GetAppDynamics wll return the app dynamics object in the ob, Else, error if nothing exists.
	*/
	GetAppDynamics() (*AppDynamics, error)

	addAppDynamicsEndpoint(endpoint string) error
	addMqEndpoint(name, application string, metricPath string) error
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}

}

type AppDynamics struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Endpoint    string
	MqEndpoints []MqEndpoint
}

/*
MqEndpoint Object
*/
type MqEndpoint struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string
	Application string
	MetricPath  string
}

func (s *mongoStore) addAppDynamicsEndpoint(endpoint string) error {
	a := AppDynamics{Endpoint: endpoint}
	b, err := s.GetAppDynamics()
	c := s.mongo.C("appDynamics")

	if err == nil {
		a.ID = b.ID
		a.MqEndpoints = b.MqEndpoints
		err = c.UpdateId(a.ID, a)
		if err != nil {
			return err
		}
	} else {
		c.Insert(a)
	}
	return nil
}

func (s *mongoStore) addMqEndpoint(name, application string, metricPath string) error {
	var mq = MqEndpoint{Application: application, MetricPath: metricPath, Name: name}
	appd, err := s.GetAppDynamics()
	if err != nil {
		return err
	}

	appd.MqEndpoints = append(appd.MqEndpoints, mq)
	c := s.mongo.C("appDynamics")
	err = c.UpdateId(appd.ID, appd)
	if err != nil {
		return err
	}
	return nil
}

func (s *mongoStore) GetAppDynamics() (*AppDynamics, error) {
	c := s.mongo.C("appDynamics")
	i, err := c.Count()
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, errors.New("no app dynamics config set")
	}
	a := AppDynamics{}
	err = c.Find(nil).One(&a)
	return &a, err
}
