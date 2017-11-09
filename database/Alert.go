package database

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
)

/**
Mongo Object that stores the group to send alert messages too.
 */
type AlertDB struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	GroupID int64
	HeartbeatGroupID int64
}

/*
AlertGroup returns the alert group
 */
func AlertGroup() (groupId int64, err error) {
	c := database.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("no alert group has been set")
	}
	result := AlertDB{}
	c.Find(nil).One(&result)
	return result.GroupID, nil
}

/*
HeartbeatGroup Returns the heartbeat group.
 */
func HeartbeatGroup() (groupId int64, err error){
	c := database.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("no heartbeat group has been set")
	}
	result := AlertDB{}
	c.Find(nil).One(&result)
	return result.HeartbeatGroupID, nil
}

/*
SetAlertGroup Sets the alert group. Overrides existing group if one already exists.
 */
func SetAlertGroup(groupId int64){
	c := database.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := AlertDB{}
		result.GroupID = groupId
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := AlertDB{}
		c.Find(nil).One(&result)
		result.GroupID = groupId
		err := c.UpdateId(result.ID,result)
		if err != nil {
			log.Panic(err)
		}
	}
}

/*
SetHeartbeatGroup Sets the alert group. Overrides existing group if one already exists.
 */
func SetHeartbeatGroup(groupId int64){
	c := database.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := AlertDB{}
		result.HeartbeatGroupID = groupId
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := AlertDB{}
		c.Find(nil).One(&result)
		result.HeartbeatGroupID = groupId
		err := c.UpdateId(result.ID,result)
		if err != nil {
			log.Panic(err)
		}
	}
}
