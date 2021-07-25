package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type NasaApodRetriever struct {
	URL               string
	responseProcessor ResponseProcessor
}

func (retriever NasaApodRetriever) Retrieve() (string, error) {
	fmt.Println("NasaApodRetriever: Entering Retrieve method")

	fmt.Println("NasaApodRetriever: Calling url")
	response, err := http.Get(retriever.URL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println("NasaApodRetriever: Reading response from url")
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	v := interface{}(nil)

	fmt.Println("NasaApodRetriever: Unmarshalling response from url")
	err = json.Unmarshal(responseBody, &v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NasaApodRetriever: Getting values from response")
	date := responseProcessor.GetVal("$.date", v)
	copyright := responseProcessor.GetVal("$.copyright", v)
	explanation := responseProcessor.GetVal("$.explanation", v)
	title := responseProcessor.GetVal("$.title", v)
	url := responseProcessor.GetVal("$.url", v)

	fmt.Println("NasaApodRetriever: Formatting quote")
	var quote = fmt.Sprintf("*%s*\nCopyright %s\n%s\n\n%s\n\n%s", title, copyright, date, explanation, url)

	fmt.Println("NasaApodRetriever: Returning quote")
	return quote, err
}
