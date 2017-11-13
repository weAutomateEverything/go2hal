package service

import (
	"github.com/go-chef/chef"
	"github.com/zamedic/go2hal/database"
)

//AddChefClient Adds a chef client.
func AddChefClient(name, key, url string) error {
	//Check if the details work
	_, err := connect(name, key, url)
	if err != nil {
		return err
	}
	//No Error - therefore we assume a successful connection
	database.AddChefClient(name, url, key)

	return nil
}

func connect(name, key, url string) (client *chef.Client, err error) {
	client, err = chef.NewClient(&chef.Config{
		Name:    name,
		Key:     key,
		BaseURL: url,
	})
	return
}
