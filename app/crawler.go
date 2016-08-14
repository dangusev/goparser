package app

import (
    "fmt"
    "github.com/dangusev/goparser/app/utils"
    "github.com/moovweb/gokogiri"
    "log"
    "net/url"
    "regexp"
    "strconv"
    "strings"
    "sync"
)

// Get price value from string
func cleanPrice(s string) int64 {
    r := regexp.MustCompile("\\w+")
    priceString := strings.Join(r.FindAllString(s, -1), "")
    priceInt, _ := strconv.ParseInt(r.FindString(strings.Replace(priceString, " ", "", -1)), 10, 64)
    return priceInt
}

func buildPageUrl(queryUrl string, pageNumber int) string {
    parsed, _ := url.Parse(queryUrl)
    q := parsed.Query()
    q.Set("p", strconv.FormatInt(int64(pageNumber), 10))
    parsed.RawQuery = q.Encode()
    return parsed.String()
}

func getLastPageNumber(pageURL string) (pageNumber int) {
    root, _ := gokogiri.ParseHtml(utils.FetchViaTOR(pageURL))
    lastPageElement, err := root.Search("//a[text()=\"Последняя\"][@class=\"pagination-page\"]")
    if err != nil {
        log.Fatal(err)
    }
    if len(lastPageElement) > 0 {
        lastPageURL := lastPageElement[0].Attribute("href").Content()
        u, _ := url.Parse(lastPageURL)
        q, _ := url.ParseQuery(u.RawQuery)
        pageNumber, _ = strconv.Atoi(q.Get("p"))
    }
    return pageNumber
}

func saveResults(results chan Advert, stop chan bool) {
    // Listen the channel "results" and save data to db
    // TODO: Save to db
    for {
        select {
        case result := <-results:
            fmt.Println(result.URL)
        case <-stop:
            return
        }
    }
}

func crawlPage(pageUrl string, results chan Advert) {
    // Crawl the page of search results. Return results urls and last page number
    parsedUrl, _ := url.Parse(pageUrl)
    root, _ := gokogiri.ParseHtml(utils.FetchViaTOR(pageUrl))
    defer root.Free()
    elements, _ := root.Search("//div[contains(@class,\"item item_table\")]/div[@class=\"description\"]")

    for _, element := range elements {
        header, _ := element.Search(element.Path() + "/h3[contains(@class, \"title\")]/a")
        about, _ := element.Search(element.Path() + "/div[@class=\"about\"]")
        advert := Advert{
            Title: strings.Trim(strings.TrimSpace(header[0].Content()), "\n"),
            URL:   fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, header[0].Attribute("href").Content()),
            Price: cleanPrice(about[0].Content()),
        }
        results <- advert
    }
}

func RunCrawler(c *utils.GlobalContext) {
    log.Println("Crawling started")
    // TODO: Get queries from db
    queries := []string{
        "https://www.avito.ru/sankt-peterburg?q=ipad+mini",
        "https://www.avito.ru/sankt-peterburg?q=ipad+mini+2",
    }

    results := make(chan Advert)
    stopSaveResults := make(chan bool)
    // Run saving in goroutine
    go saveResults(results, stopSaveResults)

    // Iterate over queries in DB
    for _, queryUrl := range queries {
        log.Printf("current url - %s", queryUrl)
        // Divide pages on groups of 10 and make requests for each page
        //loopCount := int(math.Ceil(float64(pagesCount) / 10))
        //pagesPerLoop := 10
        pagesLeftCount := getLastPageNumber(queryUrl)
        fmt.Println(pagesLeftCount)
        log.Printf("Total pages - %v", pagesLeftCount)
        for pagesLeftCount > 0 {
            // Crawling pages by 20
            requestsPerLoop := utils.MinInt(pagesLeftCount, 20)
            var wg sync.WaitGroup
            wg.Add(requestsPerLoop)

            for requestsPerLoop > 0 {
                go func(w *sync.WaitGroup, u string, c chan Advert) {
                    defer w.Done()
                    crawlPage(u, c)
                }(&wg, buildPageUrl(queryUrl, pagesLeftCount), results)
                requestsPerLoop -= 1
                pagesLeftCount -= 1
            }
            wg.Wait()
        }
    }
    // Stoping the goroutine saveResults
    stopSaveResults <- true
    log.Println("Crawling finished")
}
