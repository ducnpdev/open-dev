# Usecase

## map
- code simple use map:
```go
type Cache map[string]any

func NewCache() Cache {
	return Cache{}
}

func (ma Cache) Get(k string) (any, bool) {
	val := ma[k]
	return val, val != nil
}

func (ma Cache) Set(k string, v any) {
	ma[k] = v
}

func (ma Cache) Del(k string) {
	delete(ma, k)
}

func (ma Cache) Contains(k string) bool {
	val := ma[k]
	return val != nil
}

func (ma Cache) Keys() []string {
	keys := make([]string, 0, len(ma))
	for k := range ma {
		keys = append(keys, k)
	}
	return keys
}

```
- run:
```go
go run main.go
```
- run test:
```go
go test ./... -v --race
```
```console
=== RUN   TestMap
--- PASS: TestMap (0.00s)
PASS
ok      open-dev/usecase/map    (cached)
```


### handle race condition
- use sync.RWMutex.
```go
type Cache struct {
	data map[string]any
	sync.RWMutex
}

func NewCache() Cache {
	return Cache{
		data: make(map[string]any),
	}
}
func (ma *Cache) Get(k string) (any, bool) {
	ma.RLock()
	defer ma.RUnlock()
	val := ma.data[k]
	return val, val != nil
}
func (ma *Cache) Set(k string, v any) {
	ma.Lock()
	defer ma.Unlock()
	ma.data[k] = v
}
func (ma *Cache) Del(k string) {
	ma.Lock()
	defer ma.Unlock()
	delete(ma.data, k)
}
func (ma *Cache) Contains(k string) bool {
	ma.RLock()
	defer ma.RUnlock()
	val := ma.data[k]
	return val != nil
}
func (ma *Cache) Keys() []string {
	ma.RLock()
	defer ma.RUnlock()
	keys := make([]string, 0, len(ma.data))
	for k := range ma.data {
		keys = append(keys, k)
	}
	return keys
}
```
- run test:
```go
go test ./... -v --race
```
```console
=== RUN   TestMap
--- PASS: TestMap (0.00s)
PASS
ok      open-dev/usecase/map    1.374s
```

### handle map index

- use multi map

```go
type Cache struct {
	data map[string]any
	sync.RWMutex
}

type CacheIndex []*Cache

func NewCacheIndex(n int) CacheIndex {
	cacheIndex := make([]*Cache, n)
	for i := 0; i < n; i++ {
		cacheIndex[i] = &Cache{
			data: make(map[string]any),
		}
	}
	return cacheIndex
}

func (c CacheIndex) getIndexCache(key string) int {
	checkSum := sha1.Sum([]byte(key))
	hash := int(checkSum[0])
	return hash % len(c)
}

func (c CacheIndex) getCache(key string) *Cache {
	index := c.getIndexCache(key)
	return c[index]
}

func (c CacheIndex) Get(k string) (any, bool) {
	indexCache := c.getCache(k)
	indexCache.RLock()
	defer indexCache.RUnlock()

	val, ok := indexCache.data[k]
	return val, ok
}

func (c CacheIndex) Set(k string, v any) {
	indexCache := c.getCache(k)
	indexCache.Lock()
	defer indexCache.Unlock()
	indexCache.data[k] = v
}

func (c CacheIndex) Del(k string) {
	indexCache := c.getCache(k)
	indexCache.Lock()
	defer indexCache.Unlock()
	delete(indexCache.data, k)
}

func (c CacheIndex) Contains(k string) bool {
	indexCache := c.getCache(k)
	indexCache.RLock()
	defer indexCache.RUnlock()
	_, ok := indexCache.data[k]
	return ok
}

func (c CacheIndex) Keys() []string {
	keys := make([]string, 0)
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	wg.Add(len(c))

	for _, cacheItem := range c {
		func(m *Cache) {
			m.RLock()
			for k := range m.data {
				mutex.Lock()
				keys = append(keys, k)
				mutex.Unlock()
			}
			m.RUnlock()
			wg.Done()
		}(cacheItem)
	}

	return keys
}
```
- run test:
```go
go test ./... -v --race
```
```console
=== RUN   TestMap
--- PASS: TestMap (0.00s)
PASS
ok      open-dev/usecase/map    (cached)
```

## Resize Image:
- Demo resize image base64:
```go
package usecase

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

func ResizeImageFromBase64(imgBase64 string, newHeight int) (string, error) {
	// convert image is base64 to byte
	unbased, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return "", fmt.Errorf("cannot decode base64 err=%v", err)
	}

	r := bytes.NewReader(unbased)
	// use library imaging
	// parse reader to image
	img, err := imaging.Decode(r)
	if err != nil {
		return "", err
	}

	// calculator new width of image
	newWidth := newHeight * img.Bounds().Max.X / img.Bounds().Max.Y

	// resize new image
	nrgba := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	return toBase64(nrgba)
}

func toBase64(dst *image.NRGBA) (string, error) {
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	if err := imaging.Encode(foo, dst, imaging.JPEG); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

```
- Output:
![alt text](images/resizeImage.png)
- Link:
  - convert image to base64: https://www.base64-image.de/
  - convert base64 to image: https://codebeautify.org/base64-to-image-converter

## 