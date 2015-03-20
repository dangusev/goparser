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


func QueriesAddHandler (c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    var results []parser.Query
    responseContext := make(Context)
    session := c.GetDBSession()
    defer session.Close()
    formData := ParseJsonRequest(r)
    queries := session.DB("goparser").C("queries")

    queries.Find(bson.M{"url": formData["URL"]}).All(&results)
    if len(results) > 0 {
        w.WriteHeader(400)
        WriteError(responseContext, "URL", "Query with such URL already exists")
        renderJson(w, responseContext)
    } else {
        queries.Insert(parser.Query{URL: formData["URL"].(string), Title: formData["Title"].(string)})
        w.WriteHeader(201)
        renderJson(w, responseContext)
    }
}

func QueriesDetailHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request){
    s := c.GetDBSession()
    defer s.Close()
    q := parser.GetQueryById(s, mux.Vars(r)["id"])
    renderJson(w, Context{"query": q})
}

func QueriesUpdateHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    session := c.GetDBSession()
    defer session.Close()
    formData := ParseJsonRequest(r)

    queries := session.DB("goparser").C("queries")
    queries.Update(
        bson.M{"_id": bson.ObjectIdHex(mux.Vars(r)["id"])},
        bson.M{"url": formData["URL"], "title": formData["Title"]},
    )
    w.WriteHeader(200)
    renderJson(w, Context{})
}

func QueriesDeleteHandler(c *GlobalContext, w http.ResponseWriter, r *http.Request) {
    session := c.GetDBSession()
    defer session.Close()
    queries := session.DB("goparser").C("queries")
    queries.RemoveId(bson.ObjectIdHex(mux.Vars(r)["id"]))
    w.WriteHeader(204)
    renderJson(w, Context{})
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