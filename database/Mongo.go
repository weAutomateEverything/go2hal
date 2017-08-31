package database

import (
	"gopkg.in/mgo.v2"
	"log"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")
	session, err := mgo.Dial("localhost")
	database = session.DB("hal")
	if err != nil {
		log.Panic(err)
	}
}
