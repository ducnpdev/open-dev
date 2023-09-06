package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestNewBlake2s128(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"},
		{[]byte("abc"), "bddd813c634239723171ef3fee98579b94964e3bb1cb3e427262c8c068d52319"},
		{[]byte("hello"), "324dcf027dd4a30a932c441f365a25e86b173defa4b8e58948253471b81b72cf"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewBlake2s128(tt.in)
			if hex.EncodeToString(result) != tt.out {
				t.Errorf("want %v; got %v", tt.out, hex.EncodeToString(result))
			}
		})
	}
}
