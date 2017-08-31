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
}

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

/**
Sets the alert group. Overrides existing group if one already exists.
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
		if (err != nil){
			log.Panic(err)
		}
	}
}
