package database

import (
	"gopkg.in/mgo.v2/bson"
	"log"
)

type bot struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Token string
}

func AddBot(botKey string){
	c := database.C("bots")
	botItem := bot{}
	botItem.Token = botKey
	c.Insert(botItem)
}

func ListBots() []string {
	c := database.C("bots")
	q := c.Find(nil)
	count, err := q.Count()
	if err != nil {
		log.Panic(err)
		return nil
	}
	var bots []bot
	result := make([]string, count)
	q.All(&bots)

	for i, item := range bots {
		result[i] = item.Token
	}

	return result
}
