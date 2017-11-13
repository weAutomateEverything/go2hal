package database

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
	"errors"
)

type command struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Name, Command string
}

type Key struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Username, Key string
}

/*
AddCommand adds a predefined command to a db
 */
func AddCommand(name, commandString string) error {
	c := database.C("commands")
	name = strings.ToUpper(name)
	com := command{Name: name, Command: commandString}
	return c.Insert(com)
}

/*
FindCommand returns a command
 */
func FindCommand(name string) (string, error) {
	c := database.C("commands")
	result := command{}
	err := c.Find(bson.M{"name": strings.ToUpper(name)}).One(&result)
	if err != nil {
		return "", err
	}
	return result.Command, nil
}

func AddKey(username, key string) error {
	c := database.C("keys")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return err
	}
	if count == 0 {
		r := Key{Username: username, Key: key}
		return c.Insert(r)
	} else {
		r := Key{}
		err = q.One(&r)
		if err != nil {
			return err
		}
		r.Key = key
		r.Username = username
		return c.UpdateId(r.ID, r)
	}
}

func GetKey() (*Key, error) {
	c := database.C("keys")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("no key found")
	}

	r := Key{}
	err = q.One(&r)
	if err != nil {
		return &r,nil
	}
	return &r, nil

}
