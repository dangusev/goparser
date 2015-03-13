package parser

import (
    "io/ioutil"
    "log"
    "math"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "sync"

    "github.com/moovweb/gokogiri"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "fmt"
    "regexp"
    "text/template"
    "bytes"
    "net/smtp"
)

const MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func makeRequest(url string) []byte {
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    body, err := ioutil.ReadAll(response.Body)
    defer response.Body.Close()
    return body
}

// Get price value from string
func cleanPrice(s string) int64 {
    r := regexp.MustCompile("\\w+")
    priceString := strings.Join(r.FindAllString(s, -1), "")
    priceInt, _ := strconv.ParseInt(r.FindString(strings.Replace(priceString, " ", "", -1)), 10, 64)
    return priceInt
}

func getData(u string) []Item {
    parsedUrl, _ := url.Parse(u)
    root, _ := gokogiri.ParseHtml(makeRequest(u))
    defer root.Free()

    data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")

    items := make([]Item, len(data))
    for i, item := range data {
        header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")
        about, _ := item.Search(item.Path() + "/div[@class=\"about\"]")

        item := Item{
            Title: strings.Trim(strings.TrimSpace(header[0].Content()), "\n"),
            URL:   fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, header[0].Attribute("href").Content()),
            Price: cleanPrice(about[0].Content()),
        }
        items[i] = item
    }
    return items
}

func buildURL(u string, page int) string {
    parsed, _ := url.Parse(u)
    q := parsed.Query()
    q.Set("p", strconv.FormatInt(int64(page), 10))
    parsed.RawQuery = q.Encode()
    return parsed.String()
}

func getPagesCount(pageURL string) (count int64) {
    count = 1
    root, _ := gokogiri.ParseHtml(makeRequest(pageURL))
    elem, err := root.Search("//a[text()=\"Последняя\"][@class=\"pagination__page\"]")
    if err != nil {
        log.Fatal(err)
    }
    if len(elem) > 0 {
        lastPageURL := elem[0].Attribute("href").Content()
        u, _ := url.Parse(lastPageURL)
        q, _ := url.ParseQuery(u.RawQuery)
        count, _ = strconv.ParseInt(q.Get("p"), 10, 64)
    }
    return
}

func RunParser() {
    var results []Query
    log.Println("Started parsing...")
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")

    queries.Find(bson.M{}).Select(bson.M{"url": 1}).All(&results)

    // Iterate over queries in DB
    for _, query := range results {
        log.Println(fmt.Sprintf("Parse query %s", query.URL))
        var parsedItems []Item
        // Divide pages on groups of 10 and make requests for each page
        pagesCount := getPagesCount(query.URL)
        loopCount := int(math.Ceil(float64(pagesCount) / 10))
        pagesPerLoop := 10

        for i := 1; i <= loopCount; i++ {
            var wg sync.WaitGroup

            // Get pages count for last loop
            if i == loopCount && int(math.Mod(float64(pagesCount), 10)) > 0 {
                pagesPerLoop = int(math.Mod(float64(pagesCount), 10))
            }

            // Make slice of pages
            pages := make([]int, 0, pagesPerLoop)
            base := (i - 1) * 10
            for c := 1; c <= pagesPerLoop; c++ {
                pages = append(pages, base+c)
            }
            // Run getData concurrently
            wg.Add(pagesPerLoop)

            for k := 0; k < pagesPerLoop; k++ {
                go func(w *sync.WaitGroup, u string) {
                    defer w.Done()
                    parsedItems = append(parsedItems, getData(u)...)
                }(&wg, buildURL(query.URL, pages[k]))
            }
            wg.Wait()
        }
        // Insert parsed data in DB
        query.Items = parsedItems
        for _, item := range query.Items {
            item.Is_new = !query.ItemsContains(item)
        }
        query.Update(session.Clone())
        log.Println(fmt.Sprintf("Parsing of query %s finished", query.URL))
    }
    log.Println("Parsing finished")
}

func SendNotifications(){
    auth := smtp.PlainAuth("", "dangusev92@gmail.com", "K8qetuQunuRuspb", "smtp.gmail.com")
    to := []string{"dangusev92@gmail.com"}
    context := make(map[string]interface{})
    queries := make([]Query, 0, 0)

    updatedQueries := GetQueriesWithNewItems()
    if len(updatedQueries) > 0{
        t := template.Must(template.ParseFiles("templates/email_message.html"))
        for _, q := range updatedQueries {
            queries = append(queries, q)
        }
        context["queries"] = queries
        buf := bytes.Buffer{}
        t.Execute(&buf, context)

        message := append([]byte(MIME), buf.Bytes()...)
        err := smtp.SendMail("smtp.gmail.com:587", auth, "goparser@gmail.com", to, message)
        if err != nil {
            log.Println(err)
        } else {
            log.Println(fmt.Sprintf("Email to %s sent", to))
        }
    }
}
