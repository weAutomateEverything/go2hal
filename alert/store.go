package alert

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
	"github.com/zamedic/go2hal/database"
	"gopkg.in/mgo.v2"
)

type Store interface{
	alertGroup() (groupID int64, err error)
	heartbeatGroup() (groupID int64, err error)
	nonTechnicalGroup() (groupID int64, err error)

	setAlertGroup(AlertGroupID int64)
	setHeartbeatGroup(groupID int64)
	setNonTechnicalGroup(groupID int64)
}

type mongoStore struct {
	mongo *mgo.Database
}

type alertDB struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	GroupID           int64
	HeartbeatGroupID  int64
	NonTechnicalGroup int64
}

func NewStore(mongo *mgo.Database) Store {
	return &mongoStore{
		mongo:mongo,
	}
}

func (s * mongoStore)alertGroup() (groupID int64, err error) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("no alert group has been set")
	}
	result := alertDB{}
	c.Find(nil).One(&result)
	return result.GroupID, nil
}

func (s * mongoStore)heartbeatGroup() (groupID int64, err error) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("no heartbeat group has been set")
	}
	result := alertDB{}
	c.Find(nil).One(&result)
	return result.HeartbeatGroupID, nil
}

func (s * mongoStore)nonTechnicalGroup() (groupID int64, err error) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("no non technical group has been set")
	}
	result := alertDB{}
	c.Find(nil).One(&result)
	return result.NonTechnicalGroup, nil
}


func (s * mongoStore)setAlertGroup(AlertGroupID int64) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := alertDB{}
		result.GroupID = AlertGroupID
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := alertDB{}
		c.Find(nil).One(&result)
		result.GroupID = AlertGroupID
		err := c.UpdateId(result.ID, result)
		if err != nil {
			log.Panic(err)
		}
	}
}


func (s * mongoStore)setHeartbeatGroup(groupID int64) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := alertDB{}
		result.HeartbeatGroupID = groupID
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := alertDB{}
		c.Find(nil).One(&result)
		result.HeartbeatGroupID = groupID
		err := c.UpdateId(result.ID, result)
		if err != nil {
			log.Panic(err)
		}
	}
}


func (s * mongoStore)setNonTechnicalGroup(groupID int64) {
	c := s.mongo.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := alertDB{}
		result.NonTechnicalGroup = groupID
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := alertDB{}
		c.Find(nil).One(&result)
		result.NonTechnicalGroup = groupID
		err := c.UpdateId(result.ID, result)
		if err != nil {
			log.Panic(err)
		}
	}
}
