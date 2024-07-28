package main

import "fmt"

type Person struct {
	Name string
}

// C1
func main() {
	p := Person{}
	defer func() {
		fmt.Printf("name: %s \n", p.Name)
	}()
	p.Name = "Opendev"
}

// C2
func main1() {
	p := Person{}
	defer fmt.Printf("name %s \n", p.Name)
	p.Name = "Opendev"
}
