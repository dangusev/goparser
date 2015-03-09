package parser

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "time"
)

// Query from site
type Query struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	URL   string        `bson:"url"`
	Items []Item        `bson:"items"`
    LastParsedAt time.Time `bson:"last_parsed_at"`
}

// Update Query instance in Mongodb
func (q *Query) Update(session *mgo.Session) {
	defer session.Close()
	err := session.DB("goparser").C("queries").Update(
		bson.M{"_id": q.ID},
		bson.M{"$set": bson.M{"items": q.Items, "last_parsed_at": time.Now()}},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (q Query) itemsAsMap() (map[string]Item) {
    itemsMap := make(map[string]Item)
    for _, item := range q.Items {
        itemsMap[item.URL] = item
    }
    return itemsMap
}

func (q Query) ItemsContains (item Item) bool {
    for key, _:= range q.itemsAsMap() {
        if key == item.URL {
            return true
        }
    }
    return false
}

// Item from site
type Item struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
	URL   string        `bson:"url"`
    Is_new bool         `bson:"is_new"`
}
