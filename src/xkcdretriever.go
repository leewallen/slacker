package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type XkcdRetriever struct {
	URL               string
	responseProcessor ResponseProcessor
}

func (retriever XkcdRetriever) Retrieve() (string, error) {
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

	alt := responseProcessor.GetVal("$.alt", v)
	title := responseProcessor.GetVal("$.title", v)
	img := responseProcessor.GetVal("$.img", v)
	year := responseProcessor.GetVal("$.year", v)
	month := responseProcessor.GetVal("$.month", v)
	day := responseProcessor.GetVal("$.day", v)
	num := responseProcessor.GetVal("$.num", v)

	var quote = fmt.Sprintf("%s - *%s*\n%s-%s-%s\n\nAlt Text: %s\n\n%s", num, title, year, month, day, alt, img)

	return quote, err
}
