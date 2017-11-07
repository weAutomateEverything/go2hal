package database

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"errors"
)

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
func AddSkynetEndpoint(url, username, password string) error {
	c := database.C("skynet")
	q := c.Find(nil)
	var s Skynet
	count, err := c.Count()
	if err != nil {
		log.Printf("Error retreiving skynet record count: %s", err.Error())
		return err
	}
	if count == 0 {
		s = Skynet{Username: username, Password: password, Address: url}
		err = c.Insert(s)
		if err != nil {
			log.Printf("Error creating Skynet Recrd: %s", err.Error())
		}
	} else {
		q.One(&s)
		s.Address = url
		s.Password = password
		s.Username = username
		err = c.UpdateId(s.ID, s)
		if err != nil {
			log.Printf("Error updating skynet record: %s", err.Error())
		}
	}
	return nil
}

func GetSkynetRecord() (Skynet, error) {
	c := database.C("skynet")
	var s Skynet
	count, err := c.Count()
	if err != nil {
		log.Printf("Error getting Skynet record count: %s", err.Error())
		return s, err
	}
	if count == 0 {
		err = errors.New("no skynet endpoint found. Please create one with the rest service")
		return s, err
	}
	c.Find(nil).One(&s)
	return s,nil
}
