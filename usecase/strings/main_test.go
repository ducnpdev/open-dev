package main

import "testing"

// var strs = []string{"a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11", "a123", "b11"}
var strs = []string{"b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda", "b111231sdfsdfsaf21312sdafsda"}

func BenchmarkV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concat(strs)
	}
}

func BenchmarkV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concatv2(strs)
	}
}

func BenchmarkV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concatv3(strs)
	}
}

func BenchmarkV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Concatv4(strs)
	}
}

//  go test -bench=.
// go test -bench=. -count 5 -run=^#
// go test -bench=. -benchtime=10s
// # go test -bench=. -count 5
//  go test -bench=. -benchtime=10s -benchmem