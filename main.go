package main

import (
    "net/http"
    "log"
    "html/template"
    "bitbucket.org/dangusev/goparser/models"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "path/filepath"
)

const TEMPLATE_DIR="templates"

func getTemplate(name string) (*template.Template) {
    tmpName := filepath.Join(TEMPLATE_DIR, name)
    t, err := template.ParseFiles(tmpName)
    if err != nil {
        log.Fatal(err)
    }
    return t
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    var results []models.Query

    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")

    queries.Find(bson.M{}).All(&results)
    t := getTemplate("main.html")
    t.Execute(w, map[string]interface{}{
        "queries": results,
    })
}

func main() {
    http.HandleFunc("/", mainHandler)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
