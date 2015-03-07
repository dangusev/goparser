package models

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Query from site
type Query struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	URL   string        `bson:"url"`
	Items []Item        `bson:"items"`
}

// Update Query instance in Mongodb
func (q *Query) Update(session *mgo.Session) {
	defer session.Close()
	err := session.DB("goparser").C("queries").Update(
		bson.M{"_id": q.ID},
		bson.M{"$set": bson.M{"items": q.Items}},
	)
	if err != nil {
		log.Fatal(err)
	}
}

// Item from site
type Item struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
	URL   string        `bson:"url"`
}
