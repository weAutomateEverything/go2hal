package skynet

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Store interface {
	/*
		AddSkynetEndpoint will add a skynet endpoint, or update if it already exists
	*/
	AddSkynetEndpoint(url, username, password string) error

	/*
		GetSkynetRecord will return the skynet record in the mongo DB, else throw an error if one doesnt exist.
	*/
	GetSkynetRecord() (Skynet, error)
}

type mongoStore struct {
	mongo mgo.Database
}

/*
Skynet is a storage object for skynet data
*/
type Skynet struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Address  string
	Username string
	Password string
}

/*
AddSkynetEndpoint will add a skynet endpoint, or update if it already exists
*/
func (s *mongoStore) AddSkynetEndpoint(url, username, password string) error {
	c := s.mongo.C("skynet")
	q := c.Find(nil)
	var skynet Skynet
	count, err := c.Count()
	if err != nil {
		log.Printf("Error retreiving skynet record count: %s", err.Error())
		return err
	}
	if count == 0 {
		skynet = Skynet{Username: username, Password: password, Address: url}
		err = c.Insert(s)
		if err != nil {
			log.Printf("Error creating Skynet Recrd: %s", err.Error())
		}
	} else {
		q.One(&s)
		skynet.Address = url
		skynet.Password = password
		skynet.Username = username
		err = c.UpdateId(skynet.ID, s)
		if err != nil {
			log.Printf("Error updating skynet record: %s", err.Error())
		}
	}
	return nil
}

/*
GetSkynetRecord will return the skynet record in the mongo DB, else throw an error if one doesnt exist.
*/
func (s *mongoStore) GetSkynetRecord() (Skynet, error) {
	c := s.mongo.C("skynet")
	var skynet Skynet
	count, err := c.Count()
	if err != nil {
		log.Printf("Error getting Skynet record count: %s", err.Error())
		return skynet, err
	}
	if count == 0 {
		err = errors.New("no skynet endpoint found. Please create one with the rest service")
		return skynet, err
	}
	c.Find(nil).One(&s)
	return skynet, nil
}
