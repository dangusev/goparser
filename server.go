package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	MakeRequest("http://www.avito.ru")
}

func MakeRequest(url string) {
	// bla
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Print(body)
}
