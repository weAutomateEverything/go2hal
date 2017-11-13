package database

import (
	"os"
	"strings"
	"strconv"
	"log"
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

func mongoSSL() bool {
	sslStr:= os.Getenv("MONGO_SSL")
	ssl, err := strconv.ParseBool(sslStr)
	if err != nil {
		log.Println("Invalid boolean value for MONGO_SSL environment variable. Setting false")
		return false
	}
	return ssl
}

func mongoConnectionString() string {
	return os.Getenv("MONGO")
}





