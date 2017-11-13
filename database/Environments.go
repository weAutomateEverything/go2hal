package database

import "gopkg.in/mgo.v2/bson"

type chefEnvironment struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Environment string
}

/*
AddChefEnvironment Adds a chef environment to alert on
 */
func AddChefEnvironment(environment string){
	c := database.C("chefenvironments")
	chefEnvironment := chefEnvironment{Environment:environment}
	c.Insert(chefEnvironment)
}
