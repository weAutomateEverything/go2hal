package database

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/zamedic/go2hal/service"
)

type chefClient struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name,URL,Key string
}

//AddChefClient Adds a Chef Client to the database.
func AddChefClient(name,url,key string){
	c := database.C("bots")
	chef := chefClient{Key:key,Name:name,URL:url}
	c.Insert(chef)
}

/*
IsChefConfigured will return true if a chef client is configured.
 */
func IsChefConfigured() bool {
	c := database.C("bots")
	count, err := c.Find(nil).Count()
	if err != nil {
		service.SendError(err)
		return false
	}
	return count != 0

}
