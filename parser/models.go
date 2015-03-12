package parser

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "time"
    "fmt"
)

// Query from site
type Query struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
    Title string        `bson:"title"`
	URL   string        `bson:"url"`
	Items []Item        `bson:"items"`
    LastParsedAt time.Time `bson:"last_parsed_at"`
}

func (q *Query) GetItemsUrl() (u string) {
    return fmt.Sprintf("/query/%s/items", q.ID.Hex())
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

func (q *Query) itemsAsMap() (map[string]Item) {
    itemsMap := make(map[string]Item)
    for _, item := range q.Items {
        itemsMap[item.URL] = item
    }
    return itemsMap
}

func (q *Query) ItemsContains (item Item) bool {
    for key, _:= range q.itemsAsMap() {
        if key == item.URL {
            return true
        }
    }
    return false
}

// Returns Query by ObjectIdHex
func GetQueryById(id string) (*Query) {
    var q Query
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")

    queries.FindId(bson.ObjectIdHex(id)).One(&q)
    return &q
}

func GetOrderedQueryItems(id string) []Item {
    var q Query
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")

    pipe := queries.Pipe([]bson.M{
        {"$match": bson.M{"_id": bson.ObjectIdHex(id)}},
        {"$unwind": "$items"},
        {"$sort": bson.M{"items.is_new": 1, "items.price": 1}},
        {"$group": bson.M{"_id":"$_id", "items": bson.M{"$push": "$items"}}},
    })
    pipe.One(&q)
    return q.Items
}

// Item from site
type Item struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
	URL   string        `bson:"url"`
    Is_new bool         `bson:"is_new"`
    Price int64         `bson:"price"`
}
