package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ImageInfo struct {
	Info Data
}

type Data struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Images      Links  `json:"images"`
}
type Links struct {
	Hidpi  string
	Normal string
	Teaser string
}

func main() {

	urls := "http://localhost:8080"
	sampleSearch := "Floating"
	searchValue := url.Values{}
	searchValue.Add("val", sampleSearch)
	res, err := http.PostForm(urls, searchValue)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	resStr := string(body)

	fmt.Println(resStr)
}
