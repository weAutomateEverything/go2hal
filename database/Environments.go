package database

import "gopkg.in/mgo.v2/bson"

type chefEnvironment struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Environment, FriendlyName string
}

/*
AddChefEnvironment Adds a chef environment to alert on
 */
func AddChefEnvironment(environment, friendlyName string){
	c := database.C("chefenvironments")
	chefEnvironment := chefEnvironment{Environment:environment, FriendlyName:friendlyName}
	c.Insert(chefEnvironment)
}
