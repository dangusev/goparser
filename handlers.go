package main
import (
    "net/http"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "path/filepath"
    "github.com/dangusev/goparser/utils"
)


func mainHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    c.Templates["main.html"].ExecuteTemplate(w, "base.html", utils.Context{})
}


func templatesAjaxHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    // Returns template for angular renderer
    templateName := r.URL.Query().Get("tname")
    _, fname := filepath.Split(templateName)
    t, exists := c.Templates[templateName]
    if exists {
        t.ExecuteTemplate(w, fname, utils.Context{})
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

func QueriesListHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    var results []parser.Query

    session := c.GetDBSession()
    defer session.Close()
    queries := session.DB("goparser").C("queries")
    queries.Find(bson.M{}).All(&results)
    utils.RenderJson(w, utils.Context{"queries": results})
}


func QueriesAddHandler (c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    var results []parser.Query
    responseContext := make(utils.Context)
    session := c.GetDBSession()
    defer session.Close()
    formData := utils.ParseJsonRequest(r)
    queries := session.DB("goparser").C("queries")

    queries.Find(bson.M{"url": formData["URL"]}).All(&results)
    if len(results) > 0 {
        w.WriteHeader(400)
        utils.WriteError(responseContext, "URL", "Query with such URL already exists")
        utils.RenderJson(w, responseContext)
    } else {
        queries.Insert(parser.Query{URL: formData["URL"].(string), Title: formData["Title"].(string)})
        w.WriteHeader(201)
        utils.RenderJson(w, responseContext)
    }
}

func QueriesDetailHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request){
    s := c.GetDBSession()
    defer s.Close()
    q := parser.GetQueryById(s, mux.Vars(r)["id"])
    utils.RenderJson(w, utils.Context{"query": q})
}

func QueriesUpdateHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    session := c.GetDBSession()
    defer session.Close()
    formData := utils.ParseJsonRequest(r)

    queries := session.DB("goparser").C("queries")
    queries.Update(
        bson.M{"_id": bson.ObjectIdHex(mux.Vars(r)["id"])},
        bson.M{"url": formData["URL"], "title": formData["Title"]},
    )
    w.WriteHeader(200)
    utils.RenderJson(w, utils.Context{})
}

func QueriesDeleteHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request) {
    session := c.GetDBSession()
    defer session.Close()
    queries := session.DB("goparser").C("queries")
    queries.RemoveId(bson.ObjectIdHex(mux.Vars(r)["id"]))
    w.WriteHeader(204)
    utils.RenderJson(w, utils.Context{})
}

func ItemsListHandler(c *utils.GlobalContext, w http.ResponseWriter, r *http.Request){
    s := c.GetDBSession()
    defer s.Close()
    q := parser.GetQueryById(s, mux.Vars(r)["id"])
    utils.RenderJson(w, utils.Context{
        "query": q,
        "items": parser.GetOrderedQueryItems(s, q.ID.Hex()),
    })
}