package firstCall

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

func NewMongoStore(db *mgo.Database) Store {
	return &store{
		db: db,
	}
}

type Store interface {
	getDefaultNumber(chat uint32) (string, error)
	setDefaultNumber(chat uint32, phoneNumber string) error
}

type store struct {
	db *mgo.Database
}

func (s *store) getDefaultNumber(chat uint32) (number string, err error) {
	c := s.db.C("default_callout")
	q := c.FindId(chat)

	r := defaultCallout{}

	count, err := q.Count()
	if err != nil {
		return
	}

	if count == 0 {
		err = fmt.Errorf("No default callout number found for %v", chat)
		return
	}

	err = q.One(&r)
	if err != nil {
		return
	}

	number = r.PhoneNumber
	return

}

func (s *store) setDefaultNumber(chat uint32, phoneNumber string) (err error) {
	c := s.db.C("default_callout")
	q := c.FindId(chat)

	r := defaultCallout{}

	count, err := q.Count()
	if err != nil {
		return
	}

	if count == 0 {
		r.PhoneNumber = phoneNumber
		r.ChatId = chat
		err = c.Insert(&r)
		return
	}

	err = q.One(&r)
	if err != nil {
		return
	}
	r.PhoneNumber = phoneNumber
	err = c.UpdateId(chat, &r)
	return
}

type defaultCallout struct {
	ChatId      uint32 `json:"id" bson:"_id,omitempty"`
	PhoneNumber string
}
