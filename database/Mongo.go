package database

import (
	"gopkg.in/mgo.v2"
	"log"
	"github.com/zamedic/go2hal/config"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")
	session, err := mgo.Dial(config.MongoAddress())
	database = session.DB("hal")
	if err != nil {
		log.Panic(err)
	}
}
