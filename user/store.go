package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {
	/*
		AddUser alows for a new user to be added to the database
	*/
	AddUser(employeeNumber, CalloutName, JiraName string)

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
	EmployeeNumber string `json:"employeeNumber"`
	CallOutName    string `json:"calloutName"`
	JIRAName       string `json:"jiraName"`
}

func (s *mongoStore) AddUser(employeeNumber, CalloutName, JiraName string) {
	c := s.mongo.C("Users")
	u := User{CallOutName: CalloutName, EmployeeNumber: employeeNumber, JIRAName: JiraName}
	c.Insert(u)
}

func (s *mongoStore) FindUserByCalloutName(name string) User {
	var r User
	c := s.mongo.C("Users")
	c.Find(bson.M{"calloutname": name}).One(&r)
	return r
}
