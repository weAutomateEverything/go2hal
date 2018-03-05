package machineLearning

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Store interface {
	SaveInputRecord(reqType string, date time.Time, fields map[string]interface{}) string
	SaveAction(requestId string, action string, date time.Time, fields map[string]interface{})
}

type mongoStore struct {
	db *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store {
	return &mongoStore{db}
}

type inputRecord struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Type   string
	Fields map[string]interface{}
	Date   time.Time
}

type action struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	RequestID string
	Action    string
	Fields    map[string]interface{}
	Date      time.Time
}

func (s *mongoStore) SaveInputRecord(reqType string, date time.Time, fields map[string]interface{}) string {
	r := inputRecord{Date: date, Fields: fields, Type: reqType, ID: bson.NewObjectId()}
	s.db.C("ml_input").Insert(&r)
	return r.ID.Hex()
}
func (s *mongoStore) SaveAction(requestId string, actionType string, date time.Time, fields map[string]interface{}) {
	r := action{Date: date, Fields: fields, Action: actionType, RequestID: requestId}
	s.db.C("ml_action").Insert(&r)
}
