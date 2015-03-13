package main

import (
    "log"
    "html/template"
    "path/filepath"
    "net/http"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "github.com/robfig/cron"
    "os"
)

const TEMPLATE_DIR = "templates"

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
        "items": parser.GetOrderedQueryItems(q.ID.Hex()),
    })
}

func main() {
    args := os.Args[1:]
    if len(args) > 0 && args[0] == "with_parser" {
        // Run parser by cron
        c := cron.New()
        c.AddFunc("@every 8h", parser.RunParser)
        c.AddFunc("@every 8h", parser.SendNotifications)
    }
    // Serve static
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    r := mux.NewRouter()
    r.HandleFunc("/", mainHandler).Name("main")
    r.HandleFunc("/query/{id}/items", ItemsListHandler).Name("items-list")

    http.Handle("/", r)
    log.Println("Run goparser on localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
