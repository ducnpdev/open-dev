package cryptodemo

import (
	"crypto/md5"
	"io"
)

func Md5test(str string) []byte {
	h := md5.New()
	io.WriteString(h, str)
	return h.Sum(nil)
}
