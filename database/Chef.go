package database

import "gopkg.in/mgo.v2/bson"

type chefClient struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name,URL,Key string
}

//Adds a Chef Client to the database.
func AddChefClient(name,url,key string){
	c := database.C("bots")
	chef := chefClient{Key:key,Name:name,URL:url}
	c.Insert(chef)
}
