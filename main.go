package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/blasyrkh123/goparser/models"
	"github.com/moovweb/gokogiri"
)

func getData(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	root, _ := gokogiri.ParseHtml(body)
	defer root.Free()

	data, _ := root.Search("//div[contains(@class,\"item\")][@data-type=\"1\"]/div[@class=\"description\"]")
	for _, item := range data {
		header, _ := item.Search(item.Path() + "/h3[@class=\"title\"]/a")

		item := models.Item{Title: header[0].Content(), Url: header[0].Attribute("href").Content()}

		print(item.Url)

	}
	return true
}

func main() {
	getData("https://www.avito.ru/sankt-peterburg/zapchasti_i_aksessuary/zapchasti/dlya_avtomobiley/rulevoe_upravlenie?i=1&q=308+%D0%BF%D0%B5%D0%B6%D0%BE+%D1%80%D1%83%D0%BB%D1%8C&s=1")
}

//TODO:
// Переход по страницам
// Запись результатов поиска в БД
// Работа с запросами (CRUD)
// Уведомления
