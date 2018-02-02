package ssh

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Store interface {
	addCommand(name, commandString string) error
	findCommand(name string) (string, error)
	addKey(username, key string) error
	getKey() (*key, error)
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store {
	return &mongoStore{db}
}

type command struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Name, Command string
}

type key struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Username, Key string
}

func (s *mongoStore) addCommand(name, commandString string) error {
	c := s.mongo.C("commands")
	name = strings.ToUpper(name)
	com := command{Name: name, Command: commandString}
	return c.Insert(com)
}

func (s *mongoStore) findCommand(name string) (string, error) {
	c := s.mongo.C("commands")
	result := command{}
	err := c.Find(bson.M{"name": strings.ToUpper(name)}).One(&result)
	if err != nil {
		return "", err
	}
	return result.Command, nil
}

func (s *mongoStore) addKey(username, k string) error {
	c := s.mongo.C("keys")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		r := key{Username: username, Key: k}
		return c.Insert(r)
	}
	r := key{}
	err = q.One(&r)
	if err != nil {
		return err
	}
	r.Key = k
	r.Username = username
	return c.UpdateId(r.ID, r)

}

func (s *mongoStore) getKey() (*key, error) {
	c := s.mongo.C("keys")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("no key found")
	}

	r := key{}
	err = q.One(&r)
	if err != nil {
		return &r, nil
	}
	return &r, nil

}
