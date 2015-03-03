package models

import "gopkg.in/mgo.v2/bson"

// Item from site
type Item struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson: "title"`
	Url   string        `bson: "url"`
}

// Query from site
type Query struct {
	Id  bson.ObjectId `bson:"_id,omitempty"`
	Url string        `bson:"url"`
}

// func (Query q) Map()
