package database

import (
	"gopkg.in/mgo.v2/bson"
)

/*
ChefClient contains the name, url and key for HAL to be able to connect to CHEF.
 */
type ChefClient struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Name,URL,Key string
}

//AddChefClient Adds a Chef Client to the database.
func AddChefClient(name,url,key string){
	c := database.C("chef")
	chef := ChefClient{Key:key,Name:name,URL:url}
	c.Insert(chef)
}

/*
GetChefClientDetails returns the chef client details
 */
func GetChefClientDetails() (ChefClient,error){
	c := database.C("chef")
	var client ChefClient
	err := c.Find(nil).One(&client)
	return client,err
}

/*
IsChefConfigured will return true if a chef client is configured.
 */
func IsChefConfigured() (bool, error) {
	c := database.C("chef")
	count, err := c.Find(nil).Count()
	if err != nil {
		return false, err

	}
	return count != 0, nil

}
