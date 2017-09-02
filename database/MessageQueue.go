package database

import "gopkg.in/mgo.v2/bson"

type message struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Message   string
	ChatID    int64
	MessageID int
}

type Message struct {
	Message   string
	ChatID    int64
	MessageID int
}

func AddMessageToQueue(message string, chatID int64, messageID int) {
	msg := message{Message: message, MessageID: messageID, ChatID: chatID}
	c := database.C("MessageQueue")
	c.Insert(msg)
}

func GetMessages() []Message {
	c := database.C("MessageQueue")
	var messages []message
	c.Find(nil).All(&messages)

	result := make([]Message, len(messages))

	for i, x := range messages {
		result[i].MessageID = x.MessageID
		result[i].ChatID = x.ChatID
		result[i].Message = x.Message
		c.RemoveId(x.ID)
	}

	return result
}
