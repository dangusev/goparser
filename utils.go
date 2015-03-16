package main
import (
    "net/http"
    "encoding/json"
    "html/template"
    "log"
)

type Context map[string]interface{}

type globalContext struct {
    templates map[string]*template.Template
}

type extendedHandler struct {
    *globalContext
    h func(*globalContext, http.ResponseWriter, *http.Request)
}



// Our appHandler type will now satisify http.Handler
func (eh extendedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    eh.h(eh.globalContext, w, r)
    log.Println(r.Method, r.URL)
}


func renderJson (w http.ResponseWriter, c Context) {
    jsonEncoded, err := json.Marshal(c)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonEncoded)
}