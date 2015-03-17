package main
import (
    "net/http"
    "gopkg.in/mgo.v2"
    "log"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "path/filepath"
)


func mainHandler(c *globalContext, w http.ResponseWriter, r *http.Request) {
    c.templates["main.html"].ExecuteTemplate(w, "base.html", Context{})
}


func templatesAjaxHandler(c *globalContext, w http.ResponseWriter, r *http.Request) {
    // Returns template for angular renderer
    templateName := r.URL.Query().Get("tname")
    _, fname := filepath.Split(templateName)
    t, exists := c.templates[templateName]
    if exists {
        t.ExecuteTemplate(w, fname, Context{})
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
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