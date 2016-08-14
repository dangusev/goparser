package utils

import (
    "net/http"
    "encoding/json"
    "html/template"
    "log"
    "path/filepath"
    "github.com/gorilla/mux"
    "io/ioutil"
    "fmt"
)

type Context map[string]interface{}

type GlobalContext struct {
    // Global context of the application
    Templates map[string]*template.Template
    Router    *mux.Router
    Settings
}

type Settings struct {
    // Global settings of the application
    ProjectDir          string "/home/dan/go/src/github.com/dangusev/goparser/"
    StaticDir           string `json:"static_dir"`
    TemplateDir         string `json:"template_dir"`
    TemplateAjaxDir     string `json:"template_ajax_dir"`
    TemplateBase        string `json:"template_base"`
    EmailLogin          string `json:"email_login"`
    EmailPassword       string `json:"email_password"`
    ParserRunEveryHours int64 `json:"parser_run_every_hours"`
}

func (g *GlobalContext) PrepareSettings() {
    // Parse settings.json and save it in .Settings
    var data []byte
    var settings Settings
    data, err := ioutil.ReadFile(filepath.Join(settings.ProjectDir, "settings.json"))

    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(data, &settings)
    if err != nil {
        log.Fatal(err)
    }
    g.Settings = settings
}

type ExtendedHandler struct {
    *GlobalContext
    Get    func(*GlobalContext, http.ResponseWriter, *http.Request)
    Post   func(*GlobalContext, http.ResponseWriter, *http.Request)
    Delete func(*GlobalContext, http.ResponseWriter, *http.Request)
}

func (eh ExtendedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Handle HTTP requests according their request types
    if r.Method == "GET" {
        eh.Get(eh.GlobalContext, w, r)
        log.Println(r.Method, r.URL)
    } else if r.Method == "POST" {
        eh.Post(eh.GlobalContext, w, r)
        log.Println(r.Method, r.URL)
    } else if r.Method == "DELETE" {
        eh.Delete(eh.GlobalContext, w, r)
        log.Println(r.Method, r.URL)
    } else {
        log.Println(r.URL, fmt.Sprintf("Unsupported method: %s", r.Method))
    }
}


func MinInt(x, y int) int {
    if x < y {
        return x
    } else {
        return y
    }
}