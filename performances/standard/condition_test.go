package main

import (
	"strings"
	"testing"
)

func CheckIfElseEmpty(s string) string {
	if s == "sdfsdsdfsdsdfsdsdfsdsdfsdsdfsd" {
		return s
	}
	return s
}
func CheckIfElseLen(s string) string {
	if len(s) == 30 {
		return s
	}
	return s
}

func CheckIfElseEqual(s string) string {
	if strings.EqualFold(s, "sdfsdsdfsdsdfsdsdfsdsdfsdsdfsd") {
		return s
	}
	return s
}

func BenchmarkIfElseEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		CheckIfElseEmpty("sdfsdsdfsdsdfsdsdfsdsdfsdsdfsd")
	}
}

func BenchmarkIfElseEqual(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		CheckIfElseEqual("sdfsdsdfsdsdfsdsdfsdsdfsdsdfsd")
	}
}
func BenchmarkIfElseLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// TODO: Your Code Here
		CheckIfElseLen("anc")
	}
}
