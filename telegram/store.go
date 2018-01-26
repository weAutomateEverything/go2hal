package telegram


import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
)

type Store interface {
	SetState(user int, state string, field []string) error
	listBots() []string
	getState(user int) State
}

type mongoStore struct{
	mongo *mgo.Database
}

/*
State is the current user state when flowing through a multi step process.
 */
type State struct {
	Userid int
	State string
	Field []string
}

type bot struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Token string
	Taken time.Time
	LastUpdate time.Time
	Name string
}

type HeartBeat struct {
	Name string `json:"name,omitempty"`
	LastUpdate time.Time `json:"last_updated,omitempty"`
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

func (s mongoStore)SetState(user int, state string, field []string) error{
	ss := State{Userid:user,State:state,Field:field}
	c := s.mongo.C("userstate")
	c.RemoveAll(bson.M{"userid":user})
	return c.Insert(&ss)
}


func (s mongoStore)getState(user int) State{
	c := s.mongo.C("userstate")
	var state State
	c.Find(bson.M{"userid":user}).One(&state)
	return state
}


func (s mongoStore) addBot(botKey string) error{
	c := s.mongo.C("bots")
	botItem := bot{}
	botItem.Token = botKey
	err := c.Insert(botItem)
	if err != nil {
		return err
	}
	return nil
}

func (s mongoStore)listBots() []string {
	bots,_ := s.getBots()
	result := make([]string,len(bots))
	for i, item := range bots {
		result[i] = item.Token
	}
	return result
}

func (s mongoStore)findBot(botToken string) (bot,error) {
	c := s.mongo.C("bots")
	result := bot{}
	err := c.Find(bson.M{"token":botToken}).One(&result)
	return result,err
}



func (s mongoStore)getBots() ([]bot, error) {
	c := s.mongo.C("bots")
	q := c.Find(nil)
	var bots []bot
	err := q.All(&bots)
	return bots,err
}

