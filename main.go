package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "github.com/robfig/cron"
    "os"
    "github.com/dangusev/goparser/utils"
)


func main() {
    c := &utils.GlobalContext{}
    c.PrepareSettings()
    c.PrepareTemplates()

    args := os.Args[1:]
    if len(args) > 0 && args[0] == "with_parser" {
        // Run parser by cron
        cr := cron.New()
        cr.AddFunc("@every 8h", func() {parser.RunParser(c)})
        cr.AddFunc("@every 8h", func() {parser.SendNotifications(c)})
    }

    // Serve static
    fs := http.FileServer(http.Dir(c.StaticDir))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    r := mux.NewRouter()
    c.Router = r
    r.Handle("/", utils.ExtendedHandler{GlobalContext: c, Get: mainHandler}).Name("main")
    r.Handle("/templates/", utils.ExtendedHandler{GlobalContext: c, Get: templatesAjaxHandler}).Name("templates")

    r.Handle("/api/queries/", utils.ExtendedHandler{GlobalContext: c, Get: QueriesListHandler, Post: QueriesAddHandler}).Name("queries-list")
    r.Handle("/api/queries/{id}/", utils.ExtendedHandler{GlobalContext: c, Get:QueriesDetailHandler, Post: QueriesUpdateHandler, Delete: QueriesDeleteHandler}).Name("queries-detail")
    r.Handle("/api/queries/{id}/items/", utils.ExtendedHandler{GlobalContext: c, Get: ItemsListHandler}).Name("items-list")

    http.Handle("/", r)
    log.Println("Run goparser on localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}

//TODO:
// Deploy to Digital Ocean
// Think of making Model-based handlers
// Wonder about design