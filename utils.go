package main
import (
    "net/http"
    "encoding/json"
    "html/template"
    "log"
    "path/filepath"
    "os"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
    "io/ioutil"
    "fmt"
)

type Context map[string]interface{}

func WriteError(c Context, field, message string){
    errors := make(map[string][]string)

    fieldErrors, ok := errors[field]
    if !ok {
        fieldErrors = []string{message}
    } else {
        fieldErrors = append(fieldErrors, message)
    }
    errors[field] = fieldErrors
    fmt.Println(fieldErrors, "field")
    fmt.Println(errors, "errors")
    c["errors"] = errors
}

type GlobalContext struct {
    Templates map[string]*template.Template
    Router *mux.Router
    Settings

    masterSession *mgo.Session
}

type Settings struct {
    ProjectDir string "/home/dan/go/src/github.com/dangusev/goparser/"
    StaticDir string `json:"static_dir"`
    TemplateDir string `json:"template_dir"`
    TemplateAjaxDir string `json:"template_ajax_dir"`
    TemplateBase string `json:"template_base"`
    EmailLogin string `json:"email_login"`
    EmailPassword string `json:"email_password"`
    ParserRunEveryHours int64 `json:"parser_run_every_hours"`
}

// Returns Clone of original session (masterSession)
func (g *GlobalContext) GetDBSession() *mgo.Session {
    if g.masterSession == nil {
        s, err := mgo.Dial("localhost:27017")
        if err != nil {
            log.Fatal(err)
        }
        g.masterSession = s
    }
    return g.masterSession.Clone()
}

// Parse settings.json and save it in .Settings
func (g *GlobalContext) prepareSettings() {
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

func (g *GlobalContext) prepareTemplates() {
    // custom template delimiters since the Go default delimiters clash
    // with Angular's default.
    templates := make(map[string]*template.Template)
    templateDelims := []string{"{{%", "%}}"}
    basePath := filepath.Join(g.Settings.ProjectDir, g.Settings.TemplateDir)
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

        if templateName == g.Settings.TemplateBase {
            return nil
        }
        dirName, _ := filepath.Split(templateName)
        t := template.New(g.Settings.TemplateBase)
        t.Delims(templateDelims[0], templateDelims[1])
        if dirName == "ajax/"{
            templates[templateName] = template.Must(t.ParseFiles(path))
        } else {
            templates[templateName] = template.Must(t.ParseFiles(filepath.Join(basePath, g.Settings.TemplateBase), path))
        }

        log.Printf("Processed template %s\n", templateName)
        return err
    })
    if err != nil {
        log.Fatal(err)
    }
    g.Templates = templates
}

type extendedHandler struct {
    *GlobalContext
    Get func(*GlobalContext, http.ResponseWriter, *http.Request)
    Post func(*GlobalContext, http.ResponseWriter, *http.Request)
    Delete func(*GlobalContext, http.ResponseWriter, *http.Request)
}


// Our appHandler type will now satisify http.Handler
func (eh extendedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func ParseJsonRequest(r *http.Request) map[string]interface{}{
    parsedData := make(map[string]interface{})
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&parsedData)
    if err != nil {
        log.Fatal(err)
    }
    return parsedData
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
