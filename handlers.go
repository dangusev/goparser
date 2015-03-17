package main
import (
    "net/http"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "path/filepath"
)


func mainHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    c.Templates["main.html"].ExecuteTemplate(w, "base.html", Context{})
}


func templatesAjaxHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    // Returns template for angular renderer
    templateName := r.URL.Query().Get("tname")
    _, fname := filepath.Split(templateName)
    t, exists := c.Templates[templateName]
    if exists {
        t.ExecuteTemplate(w, fname, Context{})
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func QueriesListHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    var results []parser.Query

    session := c.GetDBSession()
    defer session.Close()
    queries := session.DB("goparser").C("queries")
    queries.Find(bson.M{}).All(&results)
    renderJson(w, Context{"queries": results})
}

func ItemsListHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request){
    s := c.GetDBSession()
    defer s.Close()
    q := parser.GetQueryById(s, mux.Vars(r)["id"])
    renderJson(w, Context{
        "query": q,
        "items": parser.GetOrderedQueryItems(s, q.ID.Hex()),
    })
}