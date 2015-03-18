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

    c := &GlobalContext{}
    c.prepareSettings()
    c.prepareTemplates()

    // Routes
    r := mux.NewRouter()
    c.Router = r
    r.Handle("/", extendedHandler{GlobalContext: c, GetHandler: mainHandler}).Name("main")
    r.Handle("/templates/", extendedHandler{GlobalContext: c, GetHandler: templatesAjaxHandler}).Name("templates")

    r.Handle("/api/queries/", extendedHandler{GlobalContext: c, GetHandler: QueriesListHandler, PostHandler: QueriesAddHandler}).Name("queries-list")
    r.Handle("/api/queries/{id}/items/", extendedHandler{GlobalContext: c, GetHandler: ItemsListHandler}).Name("items-list")

    http.Handle("/", r)
    log.Println("Run goparser on localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
