package database

import (
	"gopkg.in/mgo.v2/bson"
)

type Stats struct {
	ID     bson.ObjectId `bson:"_id,omitempty"`
	Values map[string]int64
}

//IncreaseValue increases the counter on the stats db for the messages sent
func IncreaseValue(key string) {
	stat := getRecord()
	stat.Values[key]++
	saveRecord(stat)
}

//GetStats Returns the current counter as per the stats DB
func GetStats() (Stats) {
	return getRecord()
}

func getRecord() Stats {
	c := database.C("stats")
	var stat Stats
	q := c.Find(nil)
	count, _ := q.Count()
	if count == 0 {
		stat = Stats{}
		c.Insert(stat)
	} else {
		q.One(&stat)
	}

	if stat.Values == nil {
		stat.Values = make(map[string]int64)
	}
	return stat
}

func saveRecord(stat Stats) {
	c := database.C("stats")
	c.UpdateId(stat.ID, stat)
}
