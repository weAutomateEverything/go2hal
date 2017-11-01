package database

import (
	"gopkg.in/mgo.v2"
	"log"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")
	var dialinfo mgo.DialInfo
	dialinfo.Addrs = mongoServers()
	dialinfo.Database = mongoDB()
	dialinfo.Password = mongoPassword()
	dialinfo.Username= mongoUser()
	dialinfo.ReplicaSetName  = mongoReplicaSet()
	dialinfo.Source = mongoAuthSource()

	session, err := mgo.DialWithInfo(&dialinfo)
	database = session.DB(dialinfo.Database)
	if err != nil {
		log.Panic(err)
	}
}
