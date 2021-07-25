package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SwansonQuoteRetriever struct {
	URL               string
	responseProcessor ResponseProcessor
}

func (retriever SwansonQuoteRetriever) Retrieve() (string, error) {
	response, err := http.Get(retriever.URL)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	v := interface{}(nil)

	err = json.Unmarshal(responseData, &v)
	if err != nil {
		log.Fatal(err)
	}

	quote := responseProcessor.GetVal("$[0]", v)

	return quote, err
}
