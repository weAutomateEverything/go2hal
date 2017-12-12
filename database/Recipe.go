package database

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type recipe struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Recipe string
	FriendlyName string
}

/*
AddRecipe will add a recipe to the watch list for the bot
 */
func AddRecipe(recipeName, friendlyName string) error {
	c := database.C("recipes")
	recipeItem := recipe{Recipe: recipeName, FriendlyName:friendlyName}
	return c.Insert(recipeItem)

}

/*
GetRecipes returns all the configured chef recipes. 0 length if none exists or there is an error.
 */
func GetRecipes() ([]string , error) {
	c := database.C("recipes")
	q := c.Find(nil)
	var recipes []recipe
	err := q.All(&recipes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	results := make([]string,len(recipes))
	for i, recipe := range recipes {
		results[i] = recipe.Recipe
	}
	return results, nil
}
