package main

import (
	"fmt"
	"sync"
	"time"
)

type PubSub struct {
	mu       sync.Mutex
	channels map[string][]chan string
}

func NewPubSub() *PubSub {
	return &PubSub{
		channels: make(map[string][]chan string),
	}
}

func (ps *PubSub) Subscribe(topic string) <-chan string {
	ch := make(chan string)
	ps.mu.Lock()
	ps.channels[topic] = append(ps.channels[topic], ch)
	ps.mu.Unlock()
	return ch
}

func (ps *PubSub) Publish(topic, msg string) {
	ps.mu.Lock()
	for _, ch := range ps.channels[topic] {
		ch <- msg
	}
	ps.mu.Unlock()
}

func (ps *PubSub) Close(topic string) {
	ps.mu.Lock()
	for _, ch := range ps.channels[topic] {
		close(ch)
	}
	ps.mu.Unlock()
}

func main() {
	ps := NewPubSub()

	subscriber1 := ps.Subscribe("news")
	subscriber2 := ps.Subscribe("news")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for msg := range subscriber1 {
			fmt.Println("Subscriber 1 received:", msg)
		}
	}()

	go func() {
		defer wg.Done()
		for msg := range subscriber2 {
			fmt.Println("Subscriber 2 received:", msg)
		}
	}()
	time.Sleep(time.Second)

	ps.Publish("news", "Breaking News!")
	time.Sleep(time.Second)
	fmt.Println("---")
	ps.Publish("news", "Another News!")

	time.Sleep(time.Second)
	ps.Close("news")
	wg.Wait()
}
