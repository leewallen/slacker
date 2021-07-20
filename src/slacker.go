package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const SwansonUser = "Ron Swanson"
const SwansonIcon = ":ron-swanson:"
const NasaUser = "NASA Astronomy Picture of the Day"
const NasaIcon = ":nasa:"
const XkcdUser = "XKCD"
const XkcdIcon = ":xkcd:"

var slackUrl = os.Getenv("SLACK_URL")
var swansonUrl = os.Getenv("SWANSON_URL")
var swansonChannel = os.Getenv("SWANSON_CHANNEL")
var nasaUrl = os.Getenv("NASA_URL")
var nasaChannel = os.Getenv("NASA_CHANNEL")
var xkcdUrl = os.Getenv("XKCD_URL")
var xkcdChannel = os.Getenv("XKCD_CHANNEL")

type slackMessage struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

func getHealthAndReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(""))
}

func getSwansonQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Printf("getSwansonQuote : %v\n", r)

	retriever := &SwansonQuoteRetriever{swansonUrl}

	quote, err := retriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(quote))

	sendQuoteToSlack(quote, SwansonUser, swansonChannel, SwansonIcon)
}

func getNasaApod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Printf("getNasaApod : %v\n", r)

	retriever := &NasaApodRetriever{nasaUrl}

	quote, err := retriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(quote))

	sendQuoteToSlack(quote, NasaUser, nasaChannel, NasaIcon)
}

func getXkcd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Printf("getXkcd : %v\n", r)

	retriever := &XkcdRetriever{xkcdUrl}

	quote, err := retriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(quote))

	sendQuoteToSlack(quote, XkcdUser, xkcdChannel, XkcdIcon)
}

func sendQuoteToSlack(quote, user, channel, icon string) {
	slackQuote := &slackMessage{
		Channel:   channel,
		Username:  user,
		Text:      quote,
		IconEmoji: icon,
	}

	fmt.Printf("sendQuoteToSlack : %v\n", slackQuote)

	slackQuoteJson, err := json.Marshal(slackQuote)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", slackUrl, bytes.NewBuffer(slackQuoteJson))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", responseData)
}

func handleRequests() {
	r := mux.NewRouter()
	r.HandleFunc("/swanson", getSwansonQuote)
	r.HandleFunc("/nasa", getNasaApod)
	r.HandleFunc("/xkcd", getXkcd)
	r.HandleFunc("/health", getHealthAndReadiness)
	r.HandleFunc("/readiness", getHealthAndReadiness)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	fmt.Println("Starting Slacker.")
	handleRequests()
}
