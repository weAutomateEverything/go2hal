package rest

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/zamedic/go2hal/database"
)

type recipe struct {
	FriendlyName, RecipeName string
}

func addRecipe(w http.ResponseWriter, r *http.Request) {

	var recipe recipe
	err := json.NewDecoder(r.Body).Decode(&recipe)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.AddRecipe(recipe.RecipeName,recipe.FriendlyName)
	w.WriteHeader(http.StatusOK)

}
