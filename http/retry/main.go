package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func main() {
	var waitprocess chan bool
	retryClient := retryablehttp.NewClient()
	// fmt.Printf("%#v \n", retryClient)
	retryClient.RetryMax = 3                   // Maximum number of retries
	retryClient.RetryWaitMin = 1 * time.Second // Minimum wait time before retry
	retryClient.RetryWaitMax = 3 * time.Second // Maximum wait time before retry

	// Define a custom retry policy (optional) based on HTTP response status
	retryClient.CheckRetry = retryablehttp.DefaultRetryPolicy

	retryClient.HTTPClient.Timeout = time.Millisecond * 5

	httpRes, err := retryClient.Get("https://***.execute-api.ap-southeast-1.amazonaws.com/hello")

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	if httpRes.StatusCode != 200 {
		log.Fatalf("Request failed after retries: %v", err)
	}

	fmt.Println(httpRes.StatusCode)

	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	<-waitprocess
}
