package concurrency

import (
	"fmt"
	"net/http"
)

// MainErrorHandling represent main function
// call it into function main
func MainErrorHandling() {
	// handle bad
	RunErrorBad()
	// handle error good
	RunErrorGood()
}

// CallBad represent, handle error
func CallBad(done <-chan interface{}, urls ...string) <-chan *http.Response {
	resp := make(chan *http.Response)
	go func() {
		defer close(resp)
		for _, url := range urls {
			tmpResp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				continue
			}
			select {
			case <-done:
				fmt.Println("done")
				return
			case resp <- tmpResp:
			}

		}
	}()
	return resp
}

// RunErrorBad represent, run function error
func RunErrorBad() {
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for resp := range CallBad(done, urls...) {
		fmt.Printf("Response: %v\n", resp.Status)
	}
}

/*
good solution error handling
*/

// CallResult represent is a object, with error and http-response
type CallResult struct {
	Error    error
	Response *http.Response
}

// CallGood represent, handle error
func CallGood(done <-chan interface{}, urls ...string) <-chan CallResult {
	resp := make(chan CallResult)
	go func() {
		defer close(resp)
		for _, url := range urls {
			tmpResp, err := http.Get(url)
			result := CallResult{Error: err, Response: tmpResp}
			select {
			case <-done:
				fmt.Println("done")
				return
			case resp <- result:
				fmt.Println("<-result")
			}
		}
	}()
	return resp
}

func RunErrorGood() {
	done := make(chan interface{})
	defer close(done)
	// create indicator for break follow
	count := 0
	urls := []string{"a", "https://www.google.com", "b", "c", "d"}
	for resp := range CallGood(done, urls...) {
		// error return from func CallGood, and not handle internal goroutine
		// and cast error from goroutine
		if resp.Error != nil {
			count++
			if count > 1 {
				fmt.Println("breaking!, because too many errors")
				break
			}
			continue
		}
		fmt.Printf("Response http status: %v\n", resp.Response.Status)
	}

}
