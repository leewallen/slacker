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

type SwansonQuoteRetriever struct {
	URL string
}

func (retriever *SwansonQuoteRetriever) Retrieve() (string, error) {
	v := interface{}(nil)

	response, err := http.Get(retriever.URL)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(responseData, &v)
	if err != nil {
		log.Fatal(err)
	}

	quote := retriever.GetVal("$[0]", v)

	return quote, err
}

func (retriever *SwansonQuoteRetriever) GetVal(path string, v interface{}) string {
	val, err := jsonpath.Get(path, v)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return fmt.Sprint(val)
}
