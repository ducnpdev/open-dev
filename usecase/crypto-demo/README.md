# demo crypto
- use algorith basic and most commonly
## code example:
- sha1:
```go
func NewSha1(data []byte) []byte {
	s := sha1.New()
	s.Write(data)
	return s.Sum(nil)
}
```
- sha256
```go
func NewSha256(data []byte) []byte {
	s := sha256.New()
	s.Write(data)
	return s.Sum(nil)
}
```
- md5
```go
func Md5test(str string) []byte {
	h := md5.New()
	io.WriteString(h, str)
	return h.Sum(nil)
}
```
- blake
```go
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
// NewBlake2b256 ...
func NewBlake2b256(data []byte) []byte {
	hash := blake2b.Sum256(data)
	return hash[:]
}

// NewBlake2b512 ...
func NewBlake2b512(data []byte) []byte {
	hash := blake2b.Sum512(data)
	return hash[:]
}
```