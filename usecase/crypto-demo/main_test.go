package main

import (
	"encoding/hex"
	cryptodemo "open-dev/usecase/crypto-demo/blake"

	// cryptodemomd5 "open-dev/usecase/crypto-demo/md5"
	cryptodemosha "open-dev/usecase/crypto-demo/sha"

	"testing"
)

const str = "cryptodemomd5cryptod"

// func BenchmarkBlake2s128(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		hex.EncodeToString(cryptodemo.NewBlake2s128([]byte(str)))
// 	}
// }
// func BenchmarkBlake2s256(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		hex.EncodeToString(cryptodemo.NewBlake2s256([]byte(str)))
// 	}
// }
// func BenchmarkMD5(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		hex.EncodeToString(cryptodemomd5.Md5test(str))
// 	}
// }

func BenchmarkBlake2b256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hex.EncodeToString(cryptodemo.NewBlake2b256([]byte(str)))
	}
}

func BenchmarkSha256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hex.EncodeToString(cryptodemosha.NewSha256([]byte(str)))
	}
}

func BenchmarkSha1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hex.EncodeToString(cryptodemosha.NewSha1([]byte(str)))
	}
}
