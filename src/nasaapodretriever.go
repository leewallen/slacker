package main

import (
	"fmt"
	"log"
)

// NasaApodRetriever is used to get the api content, and format the response
type NasaApodRetriever struct {
	URL     string
	slackit Slackit
}

// Retrieve the response from the api endpoint and return a formatted response, or and error.
func (retriever NasaApodRetriever) Retrieve() (string, error) {
	fmt.Println("NasaApodRetriever: Entering Retrieve method")
	v, err := slackit.Get(retriever.URL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NasaApodRetriever: Getting values from response")
	date := slackit.GetVal("$.date", v)
	explanation := slackit.GetVal("$.explanation", v)
	title := slackit.GetVal("$.title", v)
	hdurl := slackit.GetVal("$.hdurl", v)

	return fmt.Sprintf("> *%s - %s*\n> \n> %s\n\n%s", date, title, explanation, hdurl), err
}
