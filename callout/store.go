package callout

import (
	"gopkg.in/mgo.v2"
	"time"
)

func NewStore(database *mgo.Database) Store {
	return &store{mongo: database}
}

type Store interface {
	AddAck(fields map[string]string, chat uint32, phone string, name string) error
	getAcks() ([]ackRecord, error)
	DeleteAck(chat uint32) error
	Bump(chat uint32) error
}

type store struct {
	mongo *mgo.Database
}

func (s store) Bump(chat uint32) error {
	v := ackRecord{}
	err := s.mongo.C("CALLOUT_ACK").FindId(chat).One(&v)
	if err != nil {
		return err
	}
	v.Count = v.Count + 1
	v.LastSent = time.Now()
	return s.mongo.C("CALLOUT_ACK").UpdateId(chat, &v)

}

func (s store) AddAck(fields map[string]string, chat uint32, phone string, name string) error {

	c, err := s.mongo.C("CALLOUT_ACK").FindId(chat).Count()
	if err != nil {
		return err
	}

	v := ackRecord{
		Fields:   fields,
		Chat:     chat,
		Phone:    phone,
		Name:     name,
		LastSent: time.Now(),
	}
	if c == 0 {
		return s.mongo.C("CALLOUT_ACK").Insert(&v)
	}
	return s.mongo.C("CALLOUT_ACK").UpdateId(chat, &v)

}

func (s store) getAcks() (r []ackRecord, err error) {
	err = s.mongo.C("CALLOUT_ACK").Find(nil).All(&r)
	return
}

func (s store) DeleteAck(chat uint32) error {
	return s.mongo.C("CALLOUT_ACK").RemoveId(chat)
}

type ackRecord struct {
	Chat     uint32 `json:"id" bson:"_id,omitempty"`
	Fields   map[string]string
	Phone    string
	Name     string
	Count    int
	LastSent time.Time
}
