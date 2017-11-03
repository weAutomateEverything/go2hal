package database

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"log"
)

type bot struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	Token string
	Taken time.Time
	LastUpdate time.Time
	Name string
}

type HeartBeat struct {
	Name string `json:"name,omitempty"`
	LastUpdate time.Time `json:"last_updated,omitempty"`
}

func AddBot(botKey string) error{
	c := database.C("bots")
	botItem := bot{}
	botItem.Token = botKey
	err := c.Insert(botItem)
	if err != nil {
		log.Printf("error saving bot: %s",err.Error())
		return err
	}
	return nil
}

func ListBots() []string {
	bots,_ := getBots()
	result := make([]string,len(bots))
	for i, item := range bots {
		result[i] = item.Token
	}
	return result
}


func TakeBot(botKey string){
	botItem := findBot(botKey)
	botItem.Taken = time.Now()
	updateBot(botItem)
}

func HeartbeatBot(botKey, name string){
	botItem := findBot(botKey)
	botItem.LastUpdate = time.Now()
	botItem.Name = name
	updateBot(botItem)
}

func GetBotHeartbeat() []HeartBeat {
	bots,_ := getBots()
	result := make ([]HeartBeat,len(bots))
	for i,b := range bots {
		result[i] = HeartBeat{b.Name,b.LastUpdate}
	}
	return result
}

func findBot(botToken string) bot {
	c := database.C("bots")
	result := bot{}
	err := c.Find(bson.M{"token":botToken}).One(&result)
	if err != nil {
		log.Println(err)
	}
	return result
}

func updateBot(botItem bot){
	c := database.C("bots")
	err := c.UpdateId(botItem.ID,botItem)
	if err != nil {
		log.Println(err)
	}

}

func getBots() ([]bot, error) {
	c := database.C("bots")
	q := c.Find(nil)
	var bots []bot
	err := q.All(&bots)
	return bots,err
}

