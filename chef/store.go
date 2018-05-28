package chef

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Store interface {

	/*
		AddChefEnvironment Adds a chef environment to alert on
	*/
	AddChefEnvironment(environment, friendlyName string, chatid uint32) error

	/*
		GetChefEnvironments will return all the chef environments in the database
	*/
	GetChefEnvironments() ([]ChefEnvironment, error)

	/*
		GetEnvironmentForGroup will return all chef environments for a group
	*/
	GetEnvironmentForGroup(chat uint32) ([]ChefEnvironment, error)

	/*
		GetEnvironmentFromFriendlyName returns the chef environment name based on the user friendly name supplied
	*/
	GetEnvironmentFromFriendlyName(recipe string, chat uint32) (string, error)

	/*
		GetRecipesForGroup will return all the recipes configured for a group
	*/
	GetRecipesForGroup(chat uint32) ([]Recipe, error)

	/*
		AddRecipe will add a recipe to the watch list for the bot
	*/
	AddRecipe(recipeName, friendlyName string, chatid uint32) error

	/*
		GetRecipes returns all the configured chef recipes. 0 length if none exists or there is an error.
	*/
	GetRecipes() ([]Recipe, error)

	GetRecipeFromFriendlyName(recipe string, chat uint32) (string, error)
}

func NewMongoStore(mongo *mgo.Database) Store {
	return &mongoStore{mongo}
}

type mongoStore struct {
	mongo *mgo.Database
}

func (s *mongoStore) GetEnvironmentForGroup(chat uint32) ([]ChefEnvironment, error) {
	c := s.mongo.C("chefenvironments")
	q := c.Find(bson.M{"chatid": chat})
	var r []ChefEnvironment
	err := q.All(&r)
	return r, err

}

func (s *mongoStore) GetRecipesForGroup(chat uint32) ([]Recipe, error) {
	c := s.mongo.C("recipes")
	q := c.Find(bson.M{"chatid": chat})
	var r []Recipe
	err := q.All(&r)
	return r, err
}

/*
ChefEnvironment contains the chef environments this bot is allowed to use in CHEF
*/
type ChefEnvironment struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Environment  string
	FriendlyName string
	ChatID       uint32
}

/*
Recipe are the chef recipes the bot wants to interact with
*/
type Recipe struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Recipe       string
	FriendlyName string
	ChatID       uint32
}

func (s *mongoStore) AddChefEnvironment(environment, friendlyName string, chatid uint32) error {

	c := s.mongo.C("chefenvironments")
	q := c.Find(bson.M{"environment": environment, "chatid": chatid})
	n, err := q.Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("ENvironment %v already exists for group %v", environment, chatid)

	}

	chefEnvironment := ChefEnvironment{Environment: environment, FriendlyName: friendlyName, ChatID: chatid}
	return c.Insert(chefEnvironment)
}

func (s *mongoStore) GetChefEnvironments() ([]ChefEnvironment, error) {
	c := s.mongo.C("chefenvironments")
	var r []ChefEnvironment
	err := c.Find(nil).All(&r)
	return r, err
}

func (s *mongoStore) GetEnvironmentFromFriendlyName(recipe string, chat uint32) (string, error) {
	c := s.mongo.C("chefenvironments")
	var r ChefEnvironment
	err := c.Find(bson.M{"friendlyname": recipe, "chatid": chat}).One(&r)
	return r.Environment, err
}

func (s *mongoStore) AddRecipe(recipeName, friendlyName string, chatid uint32) error {
	c := s.mongo.C("recipes")
	q := c.Find(bson.M{"recipe": recipeName, "chatid": chatid})
	n, err := q.Count()
	if err != nil {
		return err
	}

	if n > 0 {
		return fmt.Errorf("recipe with name %v already exists in group %v", recipeName, chatid)
	}

	recipeItem := Recipe{Recipe: recipeName, FriendlyName: friendlyName, ChatID: chatid}
	return c.Insert(recipeItem)

}

func (s *mongoStore) GetRecipes() ([]Recipe, error) {
	c := s.mongo.C("recipes")
	q := c.Find(nil)
	var recipes []Recipe
	err := q.All(&recipes)
	if err != nil {
		return nil, err
	}
	return recipes, nil
}

func (s *mongoStore) GetRecipeFromFriendlyName(recipe string, chat uint32) (string, error) {
	c := s.mongo.C("recipes")
	var r Recipe
	err := c.Find(bson.M{"friendlyname": recipe, "chatit": chat}).One(&r)
	return r.Recipe, err
}
