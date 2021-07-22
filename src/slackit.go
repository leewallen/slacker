package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
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

var (
	slackUrl                                       = os.Getenv("SLACK_URL")
	swansonUrl                                     = os.Getenv("SWANSON_URL")
	swansonChannel                                 = os.Getenv("SWANSON_CHANNEL")
	nasaUrl                                        = os.Getenv("NASA_URL")
	nasaChannel                                    = os.Getenv("NASA_CHANNEL")
	xkcdUrl                                        = os.Getenv("XKCD_URL")
	xkcdChannel                                    = os.Getenv("XKCD_CHANNEL")
	nasaRetriever, swansonRetriever, xkcdRetriever Retriever
	responseProcessor                              Processor
)

type slackMessage struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

type Processor struct{}

type Retriever interface {
	Retrieve() (string, error)
}

type ResponseProcessor interface {
	GetVal(path string, v interface{}) string
}

func (responseProcessor Processor) GetVal(path string, v interface{}) string {
	val, err := jsonpath.Get(path, v)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return fmt.Sprint(val)

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

	if swansonRetriever == nil {
		swansonRetriever = SwansonQuoteRetriever{swansonUrl, responseProcessor}
	}

	quote, err := swansonRetriever.Retrieve()
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

	if nasaRetriever == nil {
		nasaRetriever = NasaApodRetriever{nasaUrl, responseProcessor}
	}

	quote, err := nasaRetriever.Retrieve()
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

	if xkcdRetriever == nil {
		xkcdRetriever = XkcdRetriever{nasaUrl, responseProcessor}
	}

	quote, err := xkcdRetriever.Retrieve()
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
	fmt.Println("Starting Slackit.")
	handleRequests()
}
