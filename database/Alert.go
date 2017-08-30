package database

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"log"
)

type AlertDB struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	GroupId int64
}

func AlertGroup() (groupId int64, err error) {
	c := database.C("Alert")
	count, _ := c.Count()
	if count == 0 {
		return 0, fmt.Errorf("No alert group has been set.")
	}
	result := AlertDB{}
	c.Find(nil).One(&result)
	return result.GroupId, nil
}

func SetAlertGroup(groupId int64){
	c := database.C("Alert")
	count, _ := c.Count()

	if count == 0 {
		result := AlertDB{}
		result.GroupId = groupId
		err := c.Insert(result)
		if err != nil {
			log.Panic(err)
		}
	} else {
		result := AlertDB{}
		c.Find(nil).One(&result)
		result.GroupId = groupId
		err := c.UpdateId(result.ID,result)
		if (err != nil){
			log.Panic(err)
		}
	}
}
