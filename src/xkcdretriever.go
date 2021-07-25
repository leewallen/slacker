package main

import (
	"fmt"
	"log"
)

// XkcdRetriever is used to get the api content, and format the response.
type XkcdRetriever struct {
	URL     string
	slackit Slackit
}

// Retrieve the response from the api endpoint and return a formatted response, or and error.
func (retriever XkcdRetriever) Retrieve() (string, error) {
	v, err := slackit.Get(retriever.URL)
	if err != nil {
		log.Fatal(err)
	}
	alt := slackit.GetVal("$.alt", v)
	title := slackit.GetVal("$.title", v)
	img := slackit.GetVal("$.img", v)
	year := slackit.GetVal("$.year", v)
	month := slackit.GetVal("$.month", v)
	day := slackit.GetVal("$.day", v)
	num := slackit.GetVal("$.num", v)

	return fmt.Sprintf("> *Comic %s - \"%s\"*\n> %s-%s-%s\n> \n> Alt Text: %s\n\n%s", num, title, year, month, day, alt, img), err
}
