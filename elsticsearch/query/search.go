package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func Search() {
	// Create an Elasticsearch client
	cfg := elasticsearch.Config{
		Addresses: []string{"https://e2abb97aa9b143988a13835eca6e5b8c.us-central1.gcp.cloud.es.io:443"}, // Replace with your Elasticsearch server URL
		APIKey:    "TWozQkRZb0J6QjNxMWs1djRkMk86c2t0cjFxSm1TUDJ3QS1OeVhxZnlsZw==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Define your search query
	searchQuery := `{
	  "query": {
	    "match": {
	      "unique_id5": "d0e0"
	    }
	  }
	}`

	// Create a search request
	req := esapi.SearchRequest{
		Index: []string{"search-my"}, // Replace with the name of your Elasticsearch index
		Body:  strings.NewReader(searchQuery),
	}
	now := time.Now()
	// Execute the search request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing search request: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Search request returned an error: %s", res.Status())
	}

	var searchResults struct {
		Hits struct {
			Hits []struct {
				Source Input `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	err = json.NewDecoder(res.Body).Decode(&searchResults)
	if err != nil {
		log.Fatalf("Error decoding response body: %s", err)
	}
	fmt.Println("search spent:", time.Since(now).Milliseconds())
	for _, hit := range searchResults.Hits.Hits {
		fmt.Printf("ID: %s\n", hit.Source.UUid)
		fmt.Printf("Field1: %s\n", hit.Source.UUid1)
		fmt.Printf("Field1: %s\n", hit.Source.UUid2)
		fmt.Printf("Field1: %s\n", hit.Source.UUid3)
		fmt.Printf("Field1: %s\n", hit.Source.UUid4)
		fmt.Printf("Field1: %s\n", hit.Source.UUid5)
		fmt.Printf("Field1: %s\n", hit.Source.UUid6)
		fmt.Printf("Field1: %s\n", hit.Source.UUid7)
		fmt.Printf("Field1: %s\n", hit.Source.UUid8)
	}
	// Process the search results
	fmt.Println("Search Results:")
	fmt.Println("Status Code:", res.StatusCode)
	fmt.Println("Response Body:", res.Body)
	// fmt.Println(strings.TrimSpace(string(res.Body)))
}
