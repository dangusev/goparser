package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "github.com/robfig/cron"
    "os"
)


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

    c := &globalContext{}
    c.templates = prepareTemplates(c)

    // Routes
    r := mux.NewRouter()
    r.Handle("/", extendedHandler{c, mainHandler}).Name("main")
    r.Handle("/templates/", extendedHandler{c, templatesAjaxHandler}).Name("templates")

    r.Handle("/api/queries/", extendedHandler{c, QueriesListHandler}).Name("queries-list")
    r.Handle("/api/queries/{id}/items/", extendedHandler{c, ItemsListHandler}).Name("items-list")

    http.Handle("/", r)
    log.Println("Run goparser on localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
