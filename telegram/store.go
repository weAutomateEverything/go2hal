package telegram

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	SetState(user int, chat int64, state string, field []string) error
	getState(user int, chat int64) State

	addBot(chat int64) (id uint32, err error)
	GetRoomKey(chat uint32) (roomid int64, err error)
	GetUUID(chat int64) (uuid uint32, err error)
}

type mongoStore struct {
	mongo *mgo.Database
}

/*
State is the current user state when flowing through a multi step process.
*/
type State struct {
	Userid int
	ChatId int64
	State  string
	Field  []string
}

type botRoom struct {
	ChatId int64 `bson: "_id"`
	Uuid   uint32
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore) SetState(user int, chat int64, state string, field []string) error {
	ss := State{Userid: user, State: state, Field: field, ChatId: chat}
	c := s.mongo.C("userstate")
	c.RemoveAll(bson.M{"userid": user, "chatid": chat})
	return c.Insert(&ss)
}

func (s mongoStore) getState(user int, chat int64) State {
	c := s.mongo.C("userstate")
	var state State
	c.Find(bson.M{"userid": user, "chatid": chat}).One(&state)
	return state
}

func (s mongoStore) addBot(chat int64) (id uint32, err error) {
	c := s.mongo.C("botroom")
	q := c.Find(bson.M{"chatid": chat})
	n, err := q.Count()
	if err != nil {
		return
	}

	if n > 0 {
		g := botRoom{}
		q.One(&g)
		id = g.Uuid
		return
	}

	g := botRoom{ChatId: chat, Uuid: uuid.New().ID()}
	err = c.Insert(&g)
	id = g.Uuid

	return

}

func (s mongoStore) GetRoomKey(chat uint32) (roomid int64, err error) {
	c := s.mongo.C("botroom")
	q := c.Find(bson.M{"uuid": chat})
	n, err := q.Count()
	if err != nil {
		return
	}

	if n > 0 {
		g := botRoom{}
		q.One(&g)
		roomid = g.ChatId
		return
	}
	err = fmt.Errorf("no room found for id %v", chat)
	return
}

func (s mongoStore) GetUUID(chat int64) (uuid uint32, err error) {
	c := s.mongo.C("botroom")
	q := c.Find(bson.M{"chatid": chat})
	n, err := q.Count()
	if err != nil {
		return
	}

	if n > 0 {
		g := botRoom{}
		q.One(&g)
		uuid = g.Uuid
		return
	}

	err = fmt.Errorf("no UUID could be found for chat group %v", chat)
	return
}
