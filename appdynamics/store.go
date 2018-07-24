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
	GetAppDynamics(chat uint32) (*AppDynamics, error)
	addAppDynamicsEndpoint(chat uint32, endpoint string) error
	addMqEndpoint(name, application string, metricPath string, chat uint32) error

	getAllEndpoints() ([]AppDynamics, error)
}

// swagger:model
type AppDynamics struct {
	ChatId      uint32 `json:"id" bson:"_id,omitempty"`
	Endpoint    string
	MqEndpoints []MqEndpoint
}

// swagger:model
type MqEndpoint struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string
	Application string
	MetricPath  string
	Chat        []uint32
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}

}

type mongoStore struct {
	mongo *mgo.Database
}

func (s *mongoStore) getAllEndpoints() ([]AppDynamics, error) {
	var r []AppDynamics
	c := s.mongo.C("appDynamics")

	err := c.Find(nil).All(r)

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

func (s *mongoStore) addMqEndpoint(name, application string, metricPath string, chat uint32) error {
	var mq = MqEndpoint{Application: application, MetricPath: metricPath, Name: name}
	appd, err := s.GetAppDynamics(chat)
	if err != nil {
		return err
	}
	c := s.mongo.C("appDynamics")

	for _, mq := range appd.MqEndpoints {
		if mq.Application == application && mq.MetricPath == metricPath {
			for _, id := range mq.Chat {
				if id == chat {
					return nil
				}
			}
			mq.Chat = append(mq.Chat, chat)
			return c.UpdateId(appd.ChatId, appd)
		}
	}

	appd.MqEndpoints = append(appd.MqEndpoints, mq)
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
	err = q.One(a)
	return &a, err
}
