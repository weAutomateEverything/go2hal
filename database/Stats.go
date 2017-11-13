package database

import "gopkg.in/mgo.v2/bson"

type stats struct {
	ID                       bson.ObjectId `bson:"_id,omitempty"`
	MessagesSent             int64
	AlertsReceived           int64
	AppDynamicReceived       int64
	ChefDeliveryReceived     int64
	ChefAuditMessageReceived int64
	SkynetMessageReceived    int64
}

//SendMessage increases the counter on the stats db for the messages sent
func SendMessage() {
	stat := getRecord()
	stat.MessagesSent = stat.MessagesSent + 1
	saveRecord(stat)
}

//ReceiveAlert increases the counter on the stats db for the alerts received
func ReceiveAlert() {
	stat := getRecord()
	stat.AlertsReceived = stat.AlertsReceived + 1
	saveRecord(stat)
}

//ReceiveAppynamicsMessage increases the counter on the status db for Appdynamics messages
func ReceiveAppynamicsMessage() {
	stat := getRecord()
	stat.AppDynamicReceived = stat.AppDynamicReceived + 1
	saveRecord(stat)
}

//ReceiveChefDeliveryMessage increases the counter on the stats db for Delivery messages
func ReceiveChefDeliveryMessage() {
	stat := getRecord()
	stat.ChefDeliveryReceived = stat.AppDynamicReceived + 1
	saveRecord(stat)
}

//ReceiveChefAuditMessage increased the counter on the status db for Audit messages
func ReceiveChefAuditMessage() {
	stat := getRecord()
	stat.ChefAuditMessageReceived = stat.ChefAuditMessageReceived + 1
	saveRecord(stat)
}

func ReceiveSkynetMessage() {
	stat := getRecord()
	stat.SkynetMessageReceived = stat.SkynetMessageReceived + 1
	saveRecord(stat)
}

//GetStats Returns the current counter as per the stats DB
func GetStats() (send, alerts, appdynamics, chefDelivery int64) {
	r := getRecord()
	return r.MessagesSent, r.AlertsReceived, r.AppDynamicReceived, r.ChefDeliveryReceived
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
