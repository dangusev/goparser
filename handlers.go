package main
import (
    "net/http"
    "gopkg.in/mgo.v2"
    "log"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
)


func mainHandler(c *globalContext, w http.ResponseWriter, r *http.Request) {
    c.templates["main.html"].ExecuteTemplate(w, "base.html", Context{})
}


func templatesHandler(c *globalContext, w http.ResponseWriter, r *http.Request) {
    templateName := mux.Vars(r)["tname"]
    c.templates[templateName].ExecuteTemplate(w, "base.html", Context{})
}

func QueriesListHandler(c *globalContext, w http.ResponseWriter, r *http.Request) {
    var results []parser.Query

    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")
    queries.Find(bson.M{}).All(&results)
    renderJson(w, Context{"queries": results})
}

func ItemsListHandler(c *globalContext, w http.ResponseWriter, r *http.Request){
    q := parser.GetQueryById(mux.Vars(r)["id"])
    renderJson(w, Context{
        "query": q,
        "items": parser.GetOrderedQueryItems(q.ID.Hex()),
    })
}