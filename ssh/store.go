package ssh

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type Store interface {
	addCommand(chat uint32, name, commandString string) error
	findCommand(chat uint32, name string) (string, error)
	addKey(chat uint32, username, key string) error
	getKey(chat uint32) (*key, error)
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store {
	return &mongoStore{db}
}

type command struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	ChatId  uint32        `json:"chat_id"`
	Name    string        `json:"name"`
	Command string        `json:"command"`
}

type key struct {
	ID                 bson.ObjectId `bson:"_id,omitempty"`
	ChatId             uint32        `json:"chat_id"`
	Username           string        `json:"username"`
	EncryptedBase64Key string        `json:"encrypted_base_64_key"`
}

func (s *mongoStore) addCommand(chat uint32, name, commandString string) error {
	c := s.mongo.C("commands")
	name = strings.ToUpper(name)
	com := command{
		Name:    name,
		Command: commandString,
		ChatId:  chat,
	}
	return c.Insert(com)
}

func (s *mongoStore) findCommand(chat uint32, name string) (string, error) {
	c := s.mongo.C("commands")
	result := command{}
	err := c.Find(bson.M{"name": strings.ToUpper(name), "chat_id": chat}).One(&result)
	if err != nil {
		return "", err
	}
	return result.Command, nil
}

func (s *mongoStore) addKey(chat uint32, username, baseEncrypted64Key string) error {
	c := s.mongo.C("keys")
	q := c.Find(bson.M{"chat_id": chat})
	count, err := q.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		r := key{
			ChatId:             chat,
			Username:           username,
			EncryptedBase64Key: baseEncrypted64Key}
		return c.Insert(r)
	}
	r := key{}
	err = q.One(&r)
	if err != nil {
		return err
	}
	r.EncryptedBase64Key = baseEncrypted64Key
	r.Username = username
	return c.UpdateId(r.ID, r)

}

func (s *mongoStore) getKey(chat uint32) (*key, error) {
	c := s.mongo.C("keys")
	q := c.Find(bson.M{"chat_id": chat})
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
