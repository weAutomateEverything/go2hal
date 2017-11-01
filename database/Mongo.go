package database

import (
	"gopkg.in/mgo.v2"
	"log"
	"net"
	"crypto/tls"
)

var database *mgo.Database

func init() {
	log.Println("Starting Database")
	var dialinfo mgo.DialInfo
	dialinfo.Addrs = mongoServers()
	dialinfo.Database = mongoDB()
	dialinfo.Password = mongoPassword()
	dialinfo.Username = mongoUser()
	dialinfo.ReplicaSetName = mongoReplicaSet()
	dialinfo.Source = mongoAuthSource()

	ssl := mongoSSL()

	if ssl {
		dialinfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}

	session, err := mgo.DialWithInfo(&dialinfo)
	database = session.DB(dialinfo.Database)
	if err != nil {
		log.Panic(err)
	}
}
