package database

import "gopkg.in/mgo.v2/bson"

type stats struct {
	ID             bson.ObjectId `bson:"_id,omitempty"`
	MessagesSent   int64
	AlertsReceived int64
}

func SendMessage() {
	stat := getRecord()
	stat.MessagesSent = stat.MessagesSent + 1
	saveRecord(stat)
}

func ReceiveAlert() {
	stat := getRecord()
	stat.AlertsReceived = stat.AlertsReceived + 1
	saveRecord(stat)
}

func GetStats() (send, alerts int64) {
	r := getRecord()
	return r.MessagesSent, r.AlertsReceived
}

func getRecord() stats {
	c := database.C("stats")
	var stat stats
	q := c.Find(nil)
	count, _ := q.Count()
	if count == 0 {
		stat = stats{}
		c.Insert(stat)
	} else {
		q.One(&stat)
	}
	return stat
}

func saveRecord(stat stats) {
	c := database.C("stats")
	c.UpdateId(stat.ID, stat)
}
