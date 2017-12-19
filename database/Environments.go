package database

import "gopkg.in/mgo.v2/bson"

/*
ChefEnvironment contains the chef environments this bot is allowed to use in CHEF
 */
type ChefEnvironment struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Environment  string
	FriendlyName string
}

/*
AddChefEnvironment Adds a chef environment to alert on
 */
func AddChefEnvironment(environment, friendlyName string) {
	c := database.C("chefenvironments")
	chefEnvironment := ChefEnvironment{Environment: environment, FriendlyName: friendlyName}
	c.Insert(chefEnvironment)
}

/*
GetChefEnvironments will return all the chef environments in the database
 */
func GetChefEnvironments() ([]ChefEnvironment, error) {
	c := database.C("chefenvironments")
	var r []ChefEnvironment
	err := c.Find(nil).All(&r)
	return r, err
}

/*
GetEnvironmentFromFriendlyName returns the chef environment name based on the user friendly name supplied
 */
func GetEnvironmentFromFriendlyName(recipe string) (string, error) {
	c := database.C("chefenvironments")
	var r ChefEnvironment
	err := c.Find(bson.M{"friendlyname": recipe}).One(&r)
	return r.Environment, err
}
