package main

import (
	"fmt"
	"strconv"
	"testing"
)

func MethodInt(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
func MethodItoa(i int) string {
	return strconv.Itoa(i)
}
func MethodFmt(i int) string {
	return fmt.Sprintf("%d", i)
}
func BenchmarkFmt(b *testing.B) {
	// TODO: Initialize
	number := 10
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		MethodFmt(number)
	}
}

func BenchmarkMethodInt(b *testing.B) {
	// TODO: Initialize
	number := 10
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		MethodInt(number)
	}
}
func BenchmarkMethodItoa(b *testing.B) {
	// TODO: Initialize
	number := 10
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		MethodItoa(number)
	}
}
