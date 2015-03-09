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
)

func makeRequest(url string) []byte {
    response, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    body, err := ioutil.ReadAll(response.Body)
    defer response.Body.Close()
    return body
}

func getData(url string) []Item {
    root, _ := gokogiri.ParseHtml(makeRequest(url))
    defer root.Free()

    data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")

    items := make([]Item, len(data))
    for i, item := range data {
        header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")
        item := Item{
            Title: strings.Trim(strings.TrimSpace(header[0].Content()), "\n"),
            URL:   header[0].Attribute("href").Content(),
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

    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()
    queries := session.DB("goparser").C("queries")

    queries.Find(bson.M{}).Select(bson.M{"url": 1}).All(&results)

    // Iterate over queries in DB
    for _, query := range results {
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
        query.Items = results
        for _, item := range query.Items {
            item.Is_new = query.ItemsContains(item)
        }
        query.Update(session.Clone())

    }
}

// u := "https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley?i=1&q=3"
//TODO:
// <DONE> Переход по страницам, парсинг нескольких страниц одновременно (10 за раз)
// <DONE> URL BUILDING
// <DONE>Запись результатов поиска в БД
// Прокси
// Работа с запросами (CRUD)
// Уведомления
// Логирование
// "https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley/rulevoe_upravlenie?i=1&q=308+%D0%BF%D0%B5%D0%B6%D0%BE+%D1%80%D1%83%D0%BB%D1%8C&s=1"
//
