package database

import "gopkg.in/mgo.v2/bson"

type chefClient struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name,Url,Key string
}

func AddChefClient(name,url,key string){
	c := database.C("bots")
	chef := chefClient{Key:key,Name:name,Url:url}
	c.Insert(chef)
}
