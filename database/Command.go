package database

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type command struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Name, Command string
}


/*
AddCommand adds a predefined command to a db
 */
func AddCommand(name, s string) error {
	c := database.C("commands")
	s = strings.ToUpper(s)
	com := command{Name: name, Command: s}
	return c.Insert(com)
}

/*
FindCommand returns a command
 */
func FindCommand(name string)(string, error){
	c := database.C("commands")
	result := command{}
	err := c.Find(bson.M{"Name":strings.ToUpper(name)}).One(&result)
	if err != nil {
		return "", err
	}
	return result.Command,nil
}
