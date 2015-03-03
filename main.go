package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blasyrkh123/goparser/models"
	"github.com/moovweb/gokogiri"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getData(url string, session *mgo.Session) bool {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	root, _ := gokogiri.ParseHtml(body)
	defer root.Free()
	defer session.Close()

	data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")
	for _, item := range data {
		header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")

		item := models.Item{Title: header[0].Content(), Url: header[0].Attribute("href").Content()}

		session.DB("goparser").C("items").Insert(&models.Item{Url: item.Url, Title: "test title"})
	}
	return true
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
	for _, query := range results {
		// TODO: think about goroutines
		getData(query.Url, session.Clone())
	}

	// getData("https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley/rulevoe_upravlenie?i=1&q=308+%D0%BF%D0%B5%D0%B6%D0%BE+%D1%80%D1%83%D0%BB%D1%8C&s=1")
}

//TODO:
// Переход по страницам
// Запись результатов поиска в БД
// Работа с запросами (CRUD)
// Уведомления
