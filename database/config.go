package database

import (
	"os"
	"strings"
)

func mongoServers() []string {
	s :=  os.Getenv("MONGO_SERVERS")
	return strings.Split(s,",")
}

func mongoUser() string {
return os.Getenv("MONGO_USER")
}


func mongoPassword() string {
	return os.Getenv("MONGO_PASSWORD")
}

func mongoDB() string {
	return os.Getenv("MONGO_DATABASE")

}

func mongoReplicaSet() string {
	return os.Getenv("MONGO_REPLICA_SET")
}

func mongoAuthSource() string {
	return os.Getenv("MONGO_AUTH_SOURCE")
}





