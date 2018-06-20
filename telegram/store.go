package telegram

import (
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
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
	ChatId   int64 `json:"id" bson:"_id,omitempty"`
	Uuid     uint32
	RoomName string `json:"room_name"`
}

type authRequest struct {
	ID             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	MessageID      int
	ChatId         int64
	Name           string
	Requested      time.Time
	ApprovedByName string
	ApprovedById   int
	ApprovedTime   time.Time
	Used           bool
}

type Store interface {
	SetState(user int, chat int64, state string, field []string) error
	getState(user int, chat int64) State

	addBot(chat int64, name string) (id uint32, err error)
	GetRoomKey(chat uint32) (roomid int64, err error)
	GetUUID(chat int64, name string) (uuid uint32, err error)

	newAuthRequest(msgid int, chat int64, name string) (string, error)
	approveAuthRequest(id int, chat int64, approvedByName string, approvedById int) error
	useToken(id string) (chat int64, err error)
}

type mongoStore struct {
	mongo *mgo.Database
}

func (s mongoStore) useToken(id string) (chat int64, err error) {
	c := s.mongo.C("authrequests")
	q := c.FindId(bson.ObjectIdHex(id))
	i, err := q.Count()
	if err != nil {
		return
	}

	if i == 0 {
		err = fmt.Errorf("unabkle to find approval request for id %v", id)
		return
	}
	v := authRequest{}
	q.One(&v)
	if v.Used {
		err = fmt.Errorf("approval token %v has already been used", id)
		return
	}
	if v.ApprovedById == 0 {
		err = fmt.Errorf("approval request %v has not been approved yet", id)
		return
	}
	v.Used = true
	err = c.UpdateId(v.ID, v)
	if err != nil {
		return
	}
	return v.ChatId, nil

}

func (s mongoStore) newAuthRequest(msgid int, chat int64, name string) (string, error) {
	c := s.mongo.C("authrequests")
	v := &authRequest{
		ID:        bson.NewObjectId(),
		MessageID: msgid,
		ChatId:    chat,
		Name:      name,
		Requested: time.Now(),
	}
	err := c.Insert(&v)
	if err != nil {
		return "", err
	}
	return v.ID.Hex(), nil

}

func (s mongoStore) approveAuthRequest(id int, chat int64, approvedByName string, approvedById int) error {
	c := s.mongo.C("authrequests")
	q := c.Find(bson.M{"messageid": id, "chatid": chat})
	i, err := q.Count()
	if err != nil {
		return err
	}

	if i == 0 {
		return fmt.Errorf("unabkle to find approval request for id %v", id)
	}
	v := authRequest{}
	q.One(&v)
	v.ApprovedById = approvedById
	v.ApprovedByName = approvedByName
	v.ApprovedTime = time.Now()

	return c.UpdateId(v.ID, v)

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

func (s mongoStore) addBot(chat int64, room string) (id uint32, err error) {
	c := s.mongo.C("botroom")
	q := c.FindId(chat)
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

	g := botRoom{ChatId: chat, Uuid: uuid.New().ID(), RoomName: room}
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

func (s mongoStore) GetUUID(chat int64, room string) (uuid uint32, err error) {
	c := s.mongo.C("botroom")
	q := c.FindId(chat)
	n, err := q.Count()
	if err != nil {
		return
	}

	if n > 0 {
		g := botRoom{}
		q.One(&g)
		uuid = g.Uuid

		if room != "" {
			g.RoomName = room
			c.UpdateId(chat, g)
		}
		return
	}

	err = fmt.Errorf("no UUID could be found for chat group %v", chat)
	return
}
