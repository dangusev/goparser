package main

import (
    "net/http"
    "log"
    "html/template"
    "github.com/dangusev/goparser/parser"
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
    var results []parser.Query

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
//    fs := http.FileServer(http.Dir("static"))
//    http.Handle("/static/", http.StripPrefix("/static/", fs))
//
//    http.HandleFunc("/", mainHandler)
//    err := http.ListenAndServe(":8080", nil)
//    if err != nil {
//        log.Fatal(err)
//    }
}
