package main

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	cache := NewCache()
	for i := 0; i <= 100; i++ {
		// go func(i int) {
		cache.Set(fmt.Sprint(i), i)
		// }(i)
	}
}
