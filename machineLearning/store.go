package machineLearning

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"time"
)

type Store interface {
	SaveInputRecord(reqType string, date time.Time,fields map[string]interface{}) string
	SaveAction(requestId string, action string, date time.Time, fields map[string]interface{})
}

type mongoStore struct {
	db *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store{
	return &mongoStore{db}
}
type inputRecord struct{
	ID                bson.ObjectId `bson:"_id,omitempty"`
	Type string
	Fields map[string]interface{}
	Date time.Time
}

type action struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	RequestID 	string
	Action string
	Fields map[string]interface{}
	Date time.Time
}



func (s *mongoStore) SaveInputRecord(reqType string, date time.Time,fields map[string]interface{}) string {
	r := inputRecord{Date:date,Fields:fields,Type:reqType}
	s.db.C("ml_input").Insert(&r)
	return r.ID.String()
}
func (s *mongoStore) SaveAction(requestId string, action string, date time.Time, fields map[string]interface{}){
	r := action{Date:date,Fields:fields,Action:action,RequestID:requestId}
	s.db.C("ml_action").Insert(&r)
}


