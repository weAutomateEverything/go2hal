package database

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

type AlertDB struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	groupId string
}

func AlertGroup() (groupId string, err error) {
	c := database.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return "", fmt.Errorf("No alert group has been set.")
	}
	result := AlertDB{}
	c.Find(nil).One(&result)
	return result.groupId, nil
}

func SetAlertGroup(groupId string){
	c := database.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := AlertDB{}
		result.groupId = groupId
		c.Insert(result)
	} else {
		result := AlertDB{}
		c.Find(nil).One(&result)
		result.groupId = groupId
		c.UpdateId(result.ID,result)
	}
}
