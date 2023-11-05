package cryptodemo

import (
	"crypto/sha1"
	"crypto/sha256"
)

func NewSha1(data []byte) []byte {
	s := sha1.New()
	s.Write(data)
	return s.Sum(nil)
}

func NewSha256(data []byte) []byte {
	s := sha256.New()
	s.Write(data)
	return s.Sum(nil)
}
