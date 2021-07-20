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

type NasaApodRetriever struct {
	URL string
}

func (retriever *NasaApodRetriever) Retrieve() (string, error) {
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

	date := retriever.GetVal("$.date", v)
	copyright := retriever.GetVal("$.copyright", v)
	explanation := retriever.GetVal("$.explanation", v)
	title := retriever.GetVal("$.title", v)
	url := retriever.GetVal("$.url", v)

	var quote = fmt.Sprintf("*%s*\nCopyright %s\n%s\n\n%s\n\n%s", title, copyright, date, explanation, url)

	return quote, err
}

func (retriever *NasaApodRetriever) GetVal(path string, v interface{}) string {
	val, err := jsonpath.Get(path, v)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return fmt.Sprint(val)
}
