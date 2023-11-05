package cryptodemo

import "golang.org/x/crypto/blake2s"

// NewBlake2s128 ...
func NewBlake2s128(data []byte) []byte {
	hash, _ := blake2s.New128(data)
	return hash.Sum(nil)
}

// NewBlake2s128 ...
func NewBlake2s256(data []byte) []byte {
	hash, _ := blake2s.New256(data)
	return hash.Sum(nil)
}
