package chef

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Store interface {
	/*
	GetChefClientDetails returns the chef client details
	 */
	GetChefClientDetails() (ChefClient, error)

	/*
	IsChefConfigured will return true if a chef client is configured.
	*/
	IsChefConfigured() (bool, error)

	/*
	AddChefEnvironment Adds a chef environment to alert on
	*/
	AddChefEnvironment(environment, friendlyName string)

	/*
	GetChefEnvironments will return all the chef environments in the database
	*/
	GetChefEnvironments() ([]ChefEnvironment, error)

	/*
	GetEnvironmentFromFriendlyName returns the chef environment name based on the user friendly name supplied
	*/
	GetEnvironmentFromFriendlyName(recipe string) (string, error)

	/*
	AddRecipe will add a recipe to the watch list for the bot
	*/
	AddRecipe(recipeName, friendlyName string) error

	/*
	GetRecipes returns all the configured chef recipes. 0 length if none exists or there is an error.
	*/
	GetRecipes() ([]Recipe, error)

	/*
	GetRecipeFromFriendlyName returns the chef recipe name based on the user friendly name supplied
	*/
	GetRecipeFromFriendlyName(recipe string) (string, error)

	addChefClient(name, url, key string)
}

type mongoStore struct {
	mongo *mgo.Database
}

func NewMongoStore(mongo *mgo.Database)Store{
	return &mongoStore{mongo}
}

/*
ChefClient contains the name, url and key for HAL to be able to connect to CHEF.
 */
type ChefClient struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	Name, URL, Key string
}

/*
ChefEnvironment contains the chef environments this bot is allowed to use in CHEF
 */
type ChefEnvironment struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Environment  string
	FriendlyName string
}

/*
Recipe are the chef recipes the bot wants to interact with
 */
type Recipe struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Recipe       string
	FriendlyName string
}

func (s *mongoStore) addChefClient(name, url, key string) {
	c := s.mongo.C("chef")
	chef := ChefClient{Key: key, Name: name, URL: url}
	c.Insert(chef)
}

func (s *mongoStore) GetChefClientDetails() (ChefClient, error) {
	c := s.mongo.C("chef")
	var client ChefClient
	err := c.Find(nil).One(&client)
	return client, err
}

func (s *mongoStore) IsChefConfigured() (bool, error) {
	c := s.mongo.C("chef")
	count, err := c.Find(nil).Count()
	if err != nil {
		return false, err

	}
	return count != 0, nil

}

func (s *mongoStore) AddChefEnvironment(environment, friendlyName string) {
	c := s.mongo.C("chefenvironments")
	chefEnvironment := ChefEnvironment{Environment: environment, FriendlyName: friendlyName}
	c.Insert(chefEnvironment)
}

func (s *mongoStore) GetChefEnvironments() ([]ChefEnvironment, error) {
	c := s.mongo.C("chefenvironments")
	var r []ChefEnvironment
	err := c.Find(nil).All(&r)
	return r, err
}

func (s *mongoStore) GetEnvironmentFromFriendlyName(recipe string) (string, error) {
	c := s.mongo.C("chefenvironments")
	var r ChefEnvironment
	err := c.Find(bson.M{"friendlyname": recipe}).One(&r)
	return r.Environment, err
}


func (s *mongoStore)AddRecipe(recipeName, friendlyName string) error {
	c := s.mongo.C("recipes")
	recipeItem := Recipe{Recipe: recipeName, FriendlyName: friendlyName}
	return c.Insert(recipeItem)

}


func (s *mongoStore)GetRecipes() ([]Recipe, error) {
	c := s.mongo.C("recipes")
	q := c.Find(nil)
	var recipes []Recipe
	err := q.All(&recipes)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}


func (s *mongoStore)GetRecipeFromFriendlyName(recipe string) (string, error) {
	c := s.mongo.C("recipes")
	var r Recipe
	err := c.Find(bson.M{"friendlyname": recipe}).One(&r)
	return r.Recipe, err
}
