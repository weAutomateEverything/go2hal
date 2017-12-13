package database

import (
	"gopkg.in/mgo.v2/bson"
)

type State struct {
	Userid int
	State string
	Field []string
}


func SetState(user int, state string, field []string) error{
	s := State{Userid:user,State:state,Field:field}
	c := database.C("userstate")
	c.RemoveAll(bson.M{"userid":user})
	return c.Insert(&s)
}

func GetState(user int) State{
	c := database.C("userstate")
	var s State
	c.Find(bson.M{"userid":user}).One(&s)
	return s
}


