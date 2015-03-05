package main

import (
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/blasyrkh123/goparser/models"
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

func getData(url string) []models.Item {

	root, _ := gokogiri.ParseHtml(makeRequest(url))
	defer root.Free()

	data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")

	items := make([]models.Item, len(data))
	for i, item := range data {
		header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")
		item := models.Item{
			Title: strings.Trim(strings.TrimSpace(header[0].Content()), "\n"),
			Url:   header[0].Attribute("href").Content(),
		}
		items[i] = item
	}
	return items
}

func buildURL(u string, page int64) string {
	parsed, _ := url.Parse(u)
	q := parsed.Query()
	q.Set("p", strconv.FormatInt(page, 10))
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

func main() {
	var results []models.Query
	var parsedItems []models.Item

	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	queries := session.DB("goparser").C("queries")

	queries.Find(bson.M{}).All(&results)

	// Iterate over queries in DB
	for _, query := range results {
		// Divide pages on groups of 10 and make requests for each page
		pagesCount := getPagesCount(query.Url)
		loopCount := int(math.Ceil(float64(pagesCount) / 10))

		for i := 1; i <= loopCount; i++ {
			var wg sync.WaitGroup
			if i == loopCount {
				pagesPerLoop = int(math.Dim(float64(pagesCount), 10))
			}

			wg.Add(pagesPerLoop)
			for k := 0; k < pagesPerLoop; k++ {
				go func() {
					defer wg.Done()
					parsedItems = append(parsedItems, getData(query.Url)...)
				}()
			}
			wg.Wait()

		}
	}

}

// u := "https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley?i=1&q=3"
//TODO:
// <DONE> Переход по страницам, парсинг нескольких страниц одновременно (таймаут 0.1, 10 за раз)
// URL BUILDING
// Запись результатов поиска в БД
// Работа с запросами (CRUD)
// Уведомления
// Логирование
// "https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley/rulevoe_upravlenie?i=1&q=308+%D0%BF%D0%B5%D0%B6%D0%BE+%D1%80%D1%83%D0%BB%D1%8C&s=1"
//
