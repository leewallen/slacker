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
	"time"
)

// SwansonUser is the name to display when the Ron Swanson quote slack message is sent
const SwansonUser = "Ron Swanson"

// SwansonIcon is the avatar / app icon for the Ron Swanson quote api response
const SwansonIcon = ":ron-swanson:"

// NasaUser is the name to display when the NASA APOD slack message is sent
const NasaUser = "NASA Astronomy Picture of the Day"

// NasaIcon is the avatar / app icon for the NASA api response
const NasaIcon = ":nasa:"

// XkcdUser is the name to display when the XKCD slack message is sent
const XkcdUser = "XKCD TriWeekly Knowledge Bomb"

// XkcdIcon is the avatar / app icon for the XKCD api response
const XkcdIcon = ":xkcd:"

var (
	slackURL                                       = os.Getenv("SLACK_URL")
	swansonURL                                     = os.Getenv("SWANSON_URL")
	swansonChannel                                 = os.Getenv("SWANSON_CHANNEL")
	nasaURL                                        = os.Getenv("NASA_URL")
	nasaChannel                                    = os.Getenv("NASA_CHANNEL")
	xkcdURL                                        = os.Getenv("XKCD_URL")
	xkcdChannel                                    = os.Getenv("XKCD_CHANNEL")
	apis                                           = make(map[string]api)
	targets                                        = make(map[string][]Target)
	nasaRetriever, swansonRetriever, xkcdRetriever Retriever
	slackit                                        Slackit
)

// Target contains the target channel for messages
type Target struct {
	Channel string `json:"channel"`
}

type api struct {
	URL       string `json:"url"`
	IconEmoji string `json:"icon_emoji"`
	Username  string `json:"username"`
}

type slackMessage struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

// Retriever is an interface for retrieving a formatted message from the API endpoints response data.
type Retriever interface {
	Retrieve() (string, error)
}

// Slackit reference for exposing the methods for getting an API response and extracting values from the response.
type Slackit struct{}

func init() {
	apis["NASA"] = api{URL: os.Getenv("NASA_URL"), IconEmoji: NasaIcon, Username: NasaUser}
	targets["NASA"] = make([]Target, 1)
	targets["NASA"][0] = Target{Channel: os.Getenv("NASA_CHANNEL")}

	apis["Swanson"] = api{URL: os.Getenv("SWANSON_URL"), IconEmoji: SwansonIcon, Username: SwansonUser}
	targets["Swanson"] = make([]Target, 1)
	targets["Swanson"][0] = Target{Channel: os.Getenv("SWANSON_CHANNEL")}

	apis["XKCD"] = api{URL: os.Getenv("XKCD_URL"), IconEmoji: XkcdIcon, Username: XkcdUser}
	targets["XKCD"] = make([]Target, 1)
	targets["XKCD"][0] = Target{Channel: os.Getenv("XKCD_CHANNEL")}
}

// GetVal will extract a value from JSON and return the value
func (slackit Slackit) GetVal(path string, v interface{}) string {
	val, err := jsonpath.Get(path, v)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return fmt.Sprint(val)
}

// Get will call the api endpoint and return the response or an error
func (slackit Slackit) Get(url string) (interface{}, error) {
	response, err := http.Get(url)

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

	fmt.Printf("Retrieved data: %v\n", v)

	return v, err
}

func getHealthAndReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(""))
}

func sendQuoteToCaller(quote string, w http.ResponseWriter) {
	_, err := w.Write([]byte(quote))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent quote to caller.")
}

func getSwansonQuote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if swansonRetriever == nil {
		fmt.Println("SwansonQuoteRetriever is nil - initializing it.", time.Now())
		swansonRetriever = SwansonQuoteRetriever{
			apis["Swanson"].URL,
			slackit,
		}
	}

	quote, err := swansonRetriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	sendQuoteToCaller(quote, w)
	sendQuoteToSlack(quote, apis["Swanson"], targets["Swanson"])
}

func getNasaApod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if nasaRetriever == nil {
		fmt.Println("NasaApodRetriever is nil - initializing it.", time.Now())

		nasaRetriever = NasaApodRetriever{
			apis["NASA"].URL,
			slackit,
		}
	}

	quote, err := nasaRetriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	sendQuoteToCaller(quote, w)
	sendQuoteToSlack(quote, apis["NASA"], targets["NASA"])
}

func getXkcd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if xkcdRetriever == nil {
		fmt.Println("XkcdRetriever is nil - initializing it.", time.Now())
		xkcdRetriever = XkcdRetriever{
			apis["XKCD"].URL,
			slackit,
		}
	}

	quote, err := xkcdRetriever.Retrieve()
	if err != nil {
		log.Fatal(err)
	}

	sendQuoteToCaller(quote, w)
	sendQuoteToSlack(quote, apis["XKCD"], targets["XKCD"])
}

func sendQuoteToSlack(quote string, api api, targets []Target) {
	for _, target := range targets {
		slackQuote := &slackMessage{
			Channel:   target.Channel,
			Username:  api.Username,
			Text:      quote,
			IconEmoji: api.IconEmoji,
		}

		fmt.Printf("sendQuoteToSlack : %v\n", slackQuote)

		slackQuoteJSON, err := json.Marshal(slackQuote)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(slackQuoteJSON))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Printf("Response: %s\n", responseData)
	}
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
