package database

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Recipe struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Recipe string
	FriendlyName string
}

/*
AddRecipe will add a recipe to the watch list for the bot
 */
func AddRecipe(recipeName, friendlyName string) error {
	c := database.C("recipes")
	recipeItem := Recipe{Recipe: recipeName, FriendlyName:friendlyName}
	return c.Insert(recipeItem)

}

/*
GetRecipes returns all the configured chef recipes. 0 length if none exists or there is an error.
 */
func GetRecipes() ([]Recipe , error) {
	c := database.C("recipes")
	q := c.Find(nil)
	var recipes []Recipe
	err := q.All(&recipes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return recipes, nil
}

/*
GetRecipeFromFriendlyName returns the chef recipe name based on the user friendly name supplied
 */
func GetRecipeFromFriendlyName(recipe string) (string, error){
	c := database.C("recipes")
	var r Recipe
	err := c.Find(bson.M{"friendlyname":recipe}).One(&r)
	return r.Recipe, err
}
