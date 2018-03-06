package machineLearning

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

//Store will save the Machine Learning information to a database
type Store interface {
	SaveInputRecord(reqType string, date time.Time, fields map[string]interface{}) string
	SaveAction(requestID string, action string, date time.Time, fields map[string]interface{})
}

type mongoStore struct {
	db *mgo.Database
}

//NewMongoStore will return a mongo store service to be able to store Machine Learning information to a Mongo Database
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

//SaveInputRecord will save the fields to a mongo database
//reqType - the type of request (HTTP,Telegram,File, whatever)
//date - The Date Time the action was encountered
//fields - a map[string]interface{} of any additional information, such as the request body, method, ect..
//Returns - String to correlate the input record to the output record
func (s *mongoStore) SaveInputRecord(reqType string, date time.Time, fields map[string]interface{}) string {
	r := inputRecord{Date: date, Fields: fields, Type: reqType, ID: bson.NewObjectId()}
	s.db.C("ml_input").Insert(&r)
	return r.ID.Hex()
}

//SaveAction will store the action HAL has taken in reponse to a request received.
//requestId - the ID received from the SaveInputRecord
//actionType - the action taken (Telegram, HTTP, File, ect...)
//date - the date and time the action was taken
//fields - any additional information to be saved for the action taken
func (s *mongoStore) SaveAction(requestID string, actionType string, date time.Time, fields map[string]interface{}) {
	r := action{Date: date, Fields: fields, Action: actionType, RequestID: requestID}
	s.db.C("ml_action").Insert(&r)
}
