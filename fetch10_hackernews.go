//
// Simple application for fetch Hacker News Top Stories and show its titles
//
// Top Stories API: https://hacker-news.firebaseio.com/v0/topstories?print=pretty
// Item API: https://hacker-news.firebaseio.com/v0/item/<ID>.json?print=pretty
//

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

//
// Global Variables
//
var (
	topStoriesAPIURL string
	itemAPIURL       string
)

//
// ItemData struct to keep parsed data from json
//
type ItemData struct {
	By    string
	ID    int
	Score int
	Title string
	URL   string
	Time  int
}

//
// ItemID Struct used for templating Item ID API calls
//
type ItemID struct {
	ID int
}

//
// ############################################################################
//
//                                  MAIN
//
// ############################################################################
//
func main() {

	fmt.Printf("\n###### HACKER NEWS TOP 10 STORIES #####\n\n")
	startTime := time.Now()

	//
	// Fetch all the top stories ids
	//
	topStoriesAPIURL = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

	response, _ := http.Get(topStoriesAPIURL)
	topstoriesBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	//
	// Parse the list of ids returned by the API into an array
	// and slice it for only the first 10th items
	//
	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}
	topstoriesIds := strings.FieldsFunc(string(topstoriesBody), f)
	topstoriesIds10th := topstoriesIds[:10]

	//
	// Build a template URL to receive each item ID retrieved by Top Stories API
	// and call Item API for fetch item's data
	//
	itemAPIURL = "https://hacker-news.firebaseio.com/v0/item/{{.ID}}.json?print=pretty"
	t := template.Must(template.New("URL").Parse(itemAPIURL))

	//
	// Loop through item ids then build the API URL for fetching JSON
	// data about each item
	//
	for _, itemID := range topstoriesIds10th {

		//
		// Convert itemID(string) to build ItemID
		// structure for building API URL based on this id
		//
		id, _ := strconv.Atoi(itemID)
		itemIDSource := ItemID{id}

		//
		// Apply the data in ItemID structure into
		// API URL template creating itemURLAPI string
		//
		var buffer bytes.Buffer
		t.Execute(&buffer, itemIDSource)
		itemURLAPI := buffer.String()

		//
		// Access Item API URL fetching JSON data
		//
		itemResp, err := http.Get(itemURLAPI)
		if err != nil {
			log.Fatal(err)
		}

		itemBody, _ := ioutil.ReadAll(itemResp.Body)
		itemResp.Body.Close()

		var itemJSON ItemData
		jsonErr := json.Unmarshal(itemBody, &itemJSON)
		if err != jsonErr {
			log.Fatal(jsonErr)
		}

		fmt.Printf("%v %s(%d)\n", time.Unix(int64(itemJSON.Time), 0),
			itemJSON.Title, itemJSON.Score)

		// The Hackernews API dont allow a quick access on each item
		// so it was necessary add some delay to fech the items
		time.Sleep(10 * time.Millisecond)
	}

	endTime := time.Since(startTime)
	delayTime := 10 * 10 * time.Millisecond
	fmt.Printf("\nProcess time: %v (- %v of delay time)\n\n",
		(endTime - delayTime), delayTime)

}
