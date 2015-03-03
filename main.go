package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func getData(url string, session *mgo.Session) {

	root, _ := gokogiri.ParseHtml(makeRequest(url))
	defer root.Free()
	defer session.Close()

	data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")
	for _, item := range data {
		header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")
		item := models.Item{
			Title: strings.Trim(strings.TrimSpace(header[0].Content()), "\n"),
			Url:   header[0].Attribute("href").Content(),
		}

		_, err := session.DB("goparser").C("items").Upsert(bson.M{"url": item.Url}, &models.Item{Url: item.Url, Title: item.Title})
		if err != nil {
			log.Fatal(err)
		}
	}
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
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	queries := session.DB("goparser").C("queries")

	queries.Find(bson.M{}).All(&results)
	// for _, query := range results {
	print(getPagesCount("https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley?i=1&q=3"))
	// TODO: think about goroutines
	// getData("https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley?i=1&q=3", session.Clone())
	// }

}

//TODO:
// Переход по страницам, парсинг нескольких страниц одновременно (таймаут 0.1, 10 за раз)
// Запись результатов поиска в БД
// Работа с запросами (CRUD)
// Уведомления
// "https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley/rulevoe_upravlenie?i=1&q=308+%D0%BF%D0%B5%D0%B6%D0%BE+%D1%80%D1%83%D0%BB%D1%8C&s=1"
