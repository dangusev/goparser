package main

import (
    "log"
    "html/template"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/dangusev/goparser/parser"
    "github.com/robfig/cron"
    "os"
    "path/filepath"
)

func prepareTemplates() map[string]*template.Template{
    // custom template delimiters since the Go default delimiters clash
    // with Angular's default.
    templates := make(map[string]*template.Template)
    templateDelims := []string{"{{%", "%}}"}
    basePath := "templates"
    baseTemplate := "base.html"
    // initialize the templates,
    // couldn't have used http://golang.org/pkg/html/template/#ParseGlob
    // since we have custom delimiters.
    err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // don't process folders themselves
        if info.IsDir() {
            return nil
        }
        templateName := path[len(basePath)+1:]

        if templateName == baseTemplate {
            return nil
        }

        t := template.New(baseTemplate)
        t.Delims(templateDelims[0], templateDelims[1])
        templates[templateName] = template.Must(t.ParseFiles(filepath.Join(basePath, baseTemplate), path))

        log.Printf("Processed template %s\n", templateName)
        return err
    })
    if err != nil {
        log.Fatal(err)
    }
    return templates
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

    c := &globalContext{templates:prepareTemplates()}

    // Routes
    r := mux.NewRouter()
    r.Handle("/", extendedHandler{c, mainHandler}).Name("main")
    r.Handle("/queries/", extendedHandler{c, QueriesListHandler}).Name("queries-list")
    r.Handle("/queries/{id}/items", extendedHandler{c, ItemsListHandler}).Name("items-list")

    http.Handle("/", r)
    log.Println("Run goparser on localhost:8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal(err)
    }
}
