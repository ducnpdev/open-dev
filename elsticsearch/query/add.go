package query

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
)

type Input struct {
	UUid  string `json:"unique_id"`
	UUid1 string `json:"unique_id1"`
	UUid2 string `json:"unique_id2"`
	UUid3 string `json:"unique_id3"`
	UUid4 string `json:"unique_id4"`
	UUid5 string `json:"unique_id5"`
	UUid6 string `json:"unique_id6"`
	UUid7 string `json:"unique_id7"`
	UUid8 string `json:"unique_id8"`
}

func AddDoc() {
	// currentMaxProcs := runtime.GOMAXPROCS(0)
	// fmt.Printf("Current GOMAXPROCS: %d\n", currentMaxProcs)

	newMaxProcs := 1300
	runtime.GOMAXPROCS(newMaxProcs)

	currentMaxProcs := runtime.GOMAXPROCS(0)
	fmt.Printf("Current GOMAXPROCS: %d\n", currentMaxProcs)

	cfg := elasticsearch.Config{
		Addresses: []string{"https://e2abb97aa9b143988a13835eca6e5b8c.us-central1.gcp.cloud.es.io:443"}, // Replace with your Elasticsearch server URL
		APIKey:    "TWozQkRZb0J6QjNxMWs1djRkMk86c2t0cjFxSm1TUDJ3QS1OeVhxZnlsZw==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	for i := 0; i <= 1000000; i++ {
		for j := 0; j <= 1000; j++ {
			go add(es, i)
			fmt.Println("sleep Minute")
		}
		time.Sleep(time.Second * 2)
	}

	fmt.Println("add success")
}
func add(es *elasticsearch.Client, i int) {
	now := time.Now()
	in := Input{
		UUid:  uuid.NewString(),
		UUid1: uuid.NewString(),
		UUid2: uuid.NewString(),
		UUid3: uuid.NewString(),
		UUid4: uuid.NewString(),
		UUid5: uuid.NewString(),
		UUid6: uuid.NewString(),
		UUid7: uuid.NewString(),
		UUid8: uuid.NewString(),
	}
	strS, errI := json.Marshal(in)
	if errI == nil {
		res, err := es.Index(
			"search-my",                     // Replace with your index name
			strings.NewReader(string(strS)), // Convert the JSON document to a reader
		)
		if err != nil {
			log.Fatalf("Error indexing document: %s", err)
		}
		defer res.Body.Close()
	} else {
		fmt.Println(errI)
	}
	fmt.Println("add:", i, "spent:", time.Since(now).Milliseconds())
}
