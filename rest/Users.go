package rest

import (
	"net/http"
	"encoding/json"
	"github.com/zamedic/go2hal/database"
)

type user struct {
	EmployeeNumber string
	CallOutName    string
	JIRAName       string
}

func addUser(w http.ResponseWriter, r *http.Request){
	var u user
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	database.AddUser(u.EmployeeNumber,u.CallOutName,u.JIRAName)
	w.WriteHeader(http.StatusOK)
}
