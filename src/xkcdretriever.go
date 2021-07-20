package main

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type XkcdRetriever struct {
	URL string
}

func (retriever *XkcdRetriever) Retrieve() (string, error) {
	v := interface{}(nil)
	response, err := http.Get(retriever.URL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseBody, &v)
	if err != nil {
		log.Fatal(err)
	}

	alt := retriever.GetVal("$.alt", v)
	title := retriever.GetVal("$.title", v)
	img := retriever.GetVal("$.img", v)
	year := retriever.GetVal("$.year", v)
	month := retriever.GetVal("$.month", v)
	day := retriever.GetVal("$.day", v)
	num := retriever.GetVal("$.num", v)

	var quote = fmt.Sprintf("%s - *%s*\n%s-%s-%s\n\nAlt Text: %s\n\n%s", num, title, year, month, day, alt, img)

	return quote, err
}

func (retriever *XkcdRetriever) GetVal(path string, v interface{}) string {
	val, err := jsonpath.Get(path, v)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return fmt.Sprint(val)
}
