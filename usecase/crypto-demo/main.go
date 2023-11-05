package main

import (
	"encoding/hex"
	"fmt"
	cryptodemo "open-dev/usecase/crypto-demo/blake"
	cryptodemomd5 "open-dev/usecase/crypto-demo/md5"
	cryptodemosha "open-dev/usecase/crypto-demo/sha"
)

func main() {
	str := "abc"
	fmt.Println("md5:", hex.EncodeToString(cryptodemomd5.Md5test(str)))
	fmt.Println("blake2s-128:", hex.EncodeToString(cryptodemo.NewBlake2s128([]byte(str))))
	fmt.Println("blake2s-256:", hex.EncodeToString(cryptodemo.NewBlake2s256([]byte(str))))
	fmt.Println("blake2b-256:", hex.EncodeToString(cryptodemo.NewBlake2b256([]byte(str))))

	fmt.Println("sha-1:", hex.EncodeToString(cryptodemosha.NewSha1([]byte(str))))
	fmt.Println("sha-256:", hex.EncodeToString(cryptodemosha.NewSha256([]byte(str))))
}
