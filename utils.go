package main
import (
    "net/http"
    "encoding/json"
    "html/template"
    "log"
    "path/filepath"
    "os"
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

func prepareTemplates(*globalContext) map[string]*template.Template{
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
        dirName, _ := filepath.Split(templateName)
        t := template.New(baseTemplate)
        t.Delims(templateDelims[0], templateDelims[1])
        if dirName == "ajax/"{
            templates[templateName] = template.Must(t.ParseFiles(path))
        } else {
            templates[templateName] = template.Must(t.ParseFiles(filepath.Join(basePath, baseTemplate), path))
        }

        log.Printf("Processed template %s\n", templateName)
        return err
    })
    if err != nil {
        log.Fatal(err)
    }
    return templates
}
