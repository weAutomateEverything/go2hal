package database

import "gopkg.in/mgo.v2/bson"

type User struct {
	EmployeeNumber string `json:"employeeNumber"`
	CallOutName    string `json:"calloutName"`
	JIRAName       string `json:"jiraName"`
}

func AddUser(employeeNumber, CalloutName, JiraName string) {
	c := database.C("Users")
	u := User{CallOutName: CalloutName, EmployeeNumber: employeeNumber, JIRAName: JiraName}
	c.Insert(u)
}

func FindUserByCalloutName(name string) User {
	var r User
	c := database.C("Users")
	c.Find(bson.M{"calloutname": name}).One(&r)
	return r
}