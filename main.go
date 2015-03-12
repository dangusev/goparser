package main

import (
    "net/http"
    "log"
    "html/template"
    "github.com/dangusev/goparser/parser"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "path/filepath"
    "github.com/gorilla/mux"
)

const TEMPLATE_DIR="templates"

func prepareTemplateName(name string) string {
    return filepath.Join(TEMPLATE_DIR, name)
}

func buildTemplateNames(names ...string) []string{
    var preparedNames []string
    for _, name := range names {
        preparedNames = append(preparedNames, prepareTemplateName(name))
    }
    return preparedNames
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
    t := template.Must(template.ParseFiles(buildTemplateNames("base.html", "main.html")...))
    t.Execute(w, map[string]interface{}{
        "queries": results,
    })
}

func ItemsListHandler(w http.ResponseWriter, r *http.Request){
    q := parser.GetQueryById(mux.Vars(r)["id"])
    t := template.Must(template.ParseFiles(buildTemplateNames("base.html", "items.html")...))
    t.Execute(w, map[string]interface{}{
        "query": q,
        "items": q.Items,
    })
}

func main() {
    // Serve static
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    r := mux.NewRouter()
    r.HandleFunc("/", mainHandler).Name("main")
    r.HandleFunc("/query/{id}/items", ItemsListHandler).Name("items-list")

    http.Handle("/", r)
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
