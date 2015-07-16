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
  "time"
  "text/template"
  "unicode"
)

//
// Global Variables
//
var (
  topstories_api_url string
  item_api_url       string
)

//
// Struct to keep parsed data from json
//
type ItemData struct {
  By    string
  Id    int
  Score int
  Title string
  Url   string
  Time  int
}

//
// Struct used for templating Item ID API calls
//
type ItemID struct {
  Id int
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
  start_time := time.Now()

  //
  // Fetch all the top stories ids
  //
  topstories_api_url = "https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty"

  response, _ := http.Get(topstories_api_url)
  topstories_body, _ := ioutil.ReadAll(response.Body)
  response.Body.Close()

  //
  // Parse the list of ids returned by the API into an array
  // and slice it for only the first 10th items
  //
  f := func(c rune) bool {
         return !unicode.IsLetter(c) && !unicode.IsNumber(c)
       }
  topstories_ids := strings.FieldsFunc(string(topstories_body),f)
  topstories_ids_10th := topstories_ids[:10]


  //
  // Build a template URL to receive each item ID retrieved by Top Stories API
  // and call Item API for fetch item's data 
  //
  item_api_url = "https://hacker-news.firebaseio.com/v0/item/{{.Id}}.json?print=pretty"
  t := template.Must(template.New("URL").Parse(item_api_url))

  //
  // Loop through item ids then build the API URL for fetching JSON
  // data about each item
  //
  for _, item_id := range topstories_ids_10th {

    //
    // Convert item_id(string) to build ItemID 
    // structure for building API URL based on this id
    //
    id, _ := strconv.Atoi(item_id)
    item_id_source := ItemID{id}

    //
    // Apply the data in ItemID structure into 
    // API URL template creating item_url_api string
    //
    var buffer bytes.Buffer
    t.Execute(&buffer, item_id_source)
    item_url_api := buffer.String()
    
    //
    // Access Item API URL fetching JSON data
    //
    item_resp, err := http.Get(item_url_api)
    if err != nil {
      log.Fatal(err)
    }

    item_body, _ := ioutil.ReadAll(item_resp.Body)
    item_resp.Body.Close()

    var item_json ItemData
    json_err := json.Unmarshal(item_body, &item_json)
    if err != json_err {
      log.Fatal(json_err)
    }

    fmt.Printf("%v %s(%d)\n",time.Unix(int64(item_json.Time),0), 
                             item_json.Title, item_json.Score)

    // The Hackernews API dont allow a quick access on each item
    // so it was necessary add some delay to fech the items
    time.Sleep(10 * time.Millisecond)
  }

  end_time := time.Since(start_time)
  delay_time := 10 * 10 * time.Millisecond
  fmt.Printf("\nProcess time: %v (- %v of delay time)\n\n",
              (end_time - delay_time), delay_time)

}



