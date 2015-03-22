package parser

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
    "time"
    "fmt"
    "encoding/json"
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
    return fmt.Sprintf("/queries/%s/items", q.ID.Hex())
}

func (q *Query) GetQueryUrl() (u string) {
    return fmt.Sprintf("/queries/%s/", q.ID.Hex())
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

func (q Query) MarshalJSON() ([]byte, error) {
    /* To extend existed struct we need to create new struct
       with the same fields and add extra fields what we desire*/
    return json.Marshal(struct{
        ID    bson.ObjectId `bson:"_id,omitempty"`
        Title string        `bson:"title"`
        URL   string        `bson:"url"`
        Items []Item        `bson:"items"`
        LastParsedAt time.Time `bson:"last_parsed_at"`
        ItemsURL string
        QueryURL string
    }{
        q.ID,
        q.Title,
        q.URL,
        q.Items,
        q.LastParsedAt,
        q.GetItemsUrl(),
        q.GetQueryUrl(),
    })
}


// Returns Query by ObjectIdHex
func GetQueryById(s *mgo.Session, id string) (*Query) {
    var q Query
    queries := s.DB("goparser").C("queries")
    queries.FindId(bson.ObjectIdHex(id)).One(&q)
    return &q
}

func GetOrderedQueryItems(s *mgo.Session, id string) []Item {
    var q Query
    queries := s.DB("goparser").C("queries")

    pipe := queries.Pipe([]bson.M{
        {"$match": bson.M{"_id": bson.ObjectIdHex(id)}},
        {"$unwind": "$items"},
        {"$sort": bson.M{"items.is_new": -1, "items.price": 1}},
        {"$group": bson.M{"_id":"$_id", "items": bson.M{"$push": "$items"}}},
    })
    pipe.One(&q)
    return q.Items
}

func GetQueriesWithNewItems(s *mgo.Session) []Query {
    var result []Query
    defer s.Close()
    queries := s.DB("goparser").C("queries")
    queries.Find(bson.M{"items": bson.M{"$elemMatch": bson.M{"is_new": true}}}).All(&result)
    return result
}

// Item from site
type Item struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Title string        `bson:"title"`
	URL   string        `bson:"url"`
    Is_new bool         `bson:"is_new"`
    Price int64         `bson:"price"`
}



