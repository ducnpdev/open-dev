package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Event struct {
}

func handle(e Event) {
	fmt.Println(e)
}
func main() {
	go consumer(make(<-chan Event))
	// wait
	go func() {
		http.ListenAndServe(":1234", nil)
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		<-sigs
		done <- true
	}()
	<-done
}

// bad solution:
func consumer(ch <-chan Event) {
	fmt.Println("start consumer")
	for {
		select {
		case event := <-ch:
			handle(event)
		case <-time.After(time.Nanosecond):
			log.Println("warning: no messages received")
		}
	}
}
