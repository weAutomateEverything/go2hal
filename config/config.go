package config

import (
	"os"
)

type configuration struct {
	MongoAddress string
}

var config configuration

func init() {

}

func MongoAddress() string {
	return os.Getenv("MONGO")
}

