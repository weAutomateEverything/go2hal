package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	/*
		AddUpdateUser will verify if the employee number already exists in the DB. If it does, its updated, else added
	*/
	AddUpdateUser(employeeNumber, CalloutName, JiraName string)

	/*
		FindUserByCalloutName Return a user whos details matches the callout
	*/
	FindUserByCalloutName(name string) User
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(db *mgo.Database) Store {
	return &mongoStore{db}
}

/*
User Json object
*/
type User struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	EmployeeNumber string        `json:"employeeNumber"`
	CallOutName    string        `json:"calloutName"`
	JIRAName       string        `json:"jiraName"`
}

func (s *mongoStore) AddUpdateUser(employeeNumber, CalloutName, JiraName string) {
	c := s.mongo.C("Users")
	var r User
	err := c.Find(bson.M{"employeeNumber": employeeNumber}).One(&r)
	if err == nil {
		r.CallOutName = CalloutName
		r.JIRAName = JiraName
		c.Update(r.ID, r)
	} else {
		u := User{CallOutName: CalloutName, EmployeeNumber: employeeNumber, JIRAName: JiraName}
		c.Insert(u)
	}
}

func (s *mongoStore) FindUserByCalloutName(name string) User {
	var r User
	c := s.mongo.C("Users")
	c.Find(bson.M{"calloutname": name}).One(&r)
	return r
}
