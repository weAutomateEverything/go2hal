package appdynamics

import (
	"errors"
	"gopkg.in/mgo.v2"
)

type Store interface {
	/*
		GetAppDynamics wll return the app dynamics object in the ob, Else, error if nothing exists.
	*/
	GetAppDynamics(chat uint32) (*AppDynamics, error)
	addAppDynamicsEndpoint(chat uint32, endpoint string) error
	addMqEndpoint(name, application string, metricPath string, chat uint32, ignorePrefix []string) error

	getAllEndpoints() ([]AppDynamics, error)
}

// swagger:model
type AppDynamics struct {
	ChatId      uint32 `json:"id" bson:"_id,omitempty"`
	Endpoint    string
	MqEndpoints []*MqEndpoint
}

// swagger:model
type MqEndpoint struct {
	Name          string
	Application   string
	MetricPath    string
	Chat          uint32
	MaxMessageAge float64  `json:"max_message_age" bson:"max_message_age"`
	IgnorePrefix  []string `json:"ignore_prefix"`
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}

}

type mongoStore struct {
	mongo *mgo.Database
}

func (s *mongoStore) getAllEndpoints() ([]AppDynamics, error) {
	c := s.mongo.C("appDynamics")

	q := c.Find(nil)

	count, err := q.Count()
	if err != nil {
		return nil, err
	}

	r := make([]AppDynamics, count, count)

	err = q.All(&r)

	return r, err

}

func (s *mongoStore) addAppDynamicsEndpoint(chat uint32, endpoint string) error {
	appd, err := s.GetAppDynamics(chat)
	c := s.mongo.C("appDynamics")

	//if the record already exists
	if err == nil {
		appd.Endpoint = endpoint
		return c.UpdateId(appd.ChatId, appd)
	} else {
		return c.Insert(AppDynamics{
			ChatId:   chat,
			Endpoint: endpoint,
		})
	}
}

func (s *mongoStore) addMqEndpoint(name, application string, metricPath string, chat uint32, ignorePrefix []string) error {
	var mq = MqEndpoint{Application: application, MetricPath: metricPath, Name: name, IgnorePrefix: ignorePrefix}
	appd, err := s.GetAppDynamics(chat)
	if err != nil {
		return err
	}
	c := s.mongo.C("appDynamics")

	appd.MqEndpoints = append(appd.MqEndpoints, &mq)
	return c.UpdateId(appd.ChatId, appd)

}

func (s *mongoStore) GetAppDynamics(chat uint32) (*AppDynamics, error) {
	c := s.mongo.C("appDynamics")
	q := c.FindId(chat)
	i, err := q.Count()
	if err != nil {
		return nil, err
	}
	if i == 0 {
		return nil, errors.New("no app dynamics config set")
	}
	a := AppDynamics{}
	err = q.One(&a)
	return &a, err
}
