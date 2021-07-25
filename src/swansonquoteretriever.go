package main

import (
	"log"
)

// SwansonQuoteRetriever is used to get the api content, and format the response
type SwansonQuoteRetriever struct {
	URL     string
	slackit Slackit
}

// Retrieve the response from the api endpoint and return a formatted response, or and error.
func (retriever SwansonQuoteRetriever) Retrieve() (string, error) {
	v, err := slackit.Get(retriever.URL)
	if err != nil {
		log.Fatal(err)
	}

	return slackit.GetVal("$[0]", v), err
}
