package main

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2s"
)

func main() {
	out := NewBlake2s128([]byte("abc"))
	fmt.Println(hex.EncodeToString(out))
}

// NewBlake2b128 ...
func NewBlake2s128(data []byte) []byte {
	hash, _ := blake2s.New128(data)
	return hash.Sum(nil)
}
