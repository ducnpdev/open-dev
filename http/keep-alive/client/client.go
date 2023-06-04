package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	fmt.Println("client keep alive")

	c := &http.Client{
		Transport: &http.Transport{
			// DisableKeepAlives: true,
		},
	}
	req, err := http.NewRequest("Get", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", *req)

	for i := 0; i < 3; i++ {
		resp, err := c.Do(req)
		if err != nil {
			fmt.Println("http get error:", err)
			return
		}
		defer resp.Body.Close()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("read body error:", err)
			return
		}
		fmt.Println("response body:", string(b))
		time.Sleep(time.Second)
	}

	// step2

	// c := &http.Client{}
	// req, err := http.NewRequest("Get", "http://localhost:8080", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// for i := 0; i < 5; i++ {
	// 	log.Printf("round %d begin:\n", i+1)
	// 	for j := 0; j < i+1; j++ {
	// 		resp, err := c.Do(req)
	// 		if err != nil {
	// 			fmt.Println("http get error:", err)
	// 			return
	// 		}
	// 		defer resp.Body.Close()

	// 		b, err := ioutil.ReadAll(resp.Body)
	// 		if err != nil {
	// 			fmt.Println("read body error:", err)
	// 			return
	// 		}
	// 		log.Println("response body:", string(b))
	// 	}
	// 	log.Printf("round %d end\n", i+1)
	// 	time.Sleep(7 * time.Second)
	// }
}
