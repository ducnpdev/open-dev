package main

import (
	"fmt"
	"testing"
)

type ListFake struct {
	Name  string
	Add   string
	Phone string
	Age   int
}

var interator = 100

func CreateNonArray() (list []ListFake) {
	for i := 0; i <= interator; i++ {
		list = append(list, ListFake{
			Name:  fmt.Sprintf("%s %d", "this is name ", i),
			Add:   fmt.Sprintf("%s %d", "this is address ", i),
			Phone: fmt.Sprintf("%s %d", "this is phone ", i),
			Age:   i,
		})
	}
	return list
}

func CreatePointArray() (list []*ListFake) {
	for i := 0; i <= interator; i++ {
		list = append(list, &ListFake{
			Name:  fmt.Sprintf("%s %d", "this is name ", i),
			Add:   fmt.Sprintf("%s %d", "this is address ", i),
			Phone: fmt.Sprintf("%s %d", "this is phone ", i),
			Age:   i,
		})
	}
	return list
}

func CopyArrPointer() {
	l := CreateNonArray()
	point := make([]ListFake, len(l))
	copy(point, l)
}
func CopyArrNonPointer() {
	l := CreatePointArray()
	point := make([]*ListFake, len(l))
	copy(point, l)
}

func BenchmarkNonPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyArrNonPointer()
	}
}

func BenchmarkPoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyArrPointer()
	}
}
