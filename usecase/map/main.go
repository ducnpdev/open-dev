package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"sync"
)

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

func Run() {

	cache := NewCacheIndex(10)

	fmt.Printf("index: %d", cache.getIndexCache("1"))

	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)

	keys := cache.Keys()

	for key := range keys {
		log.Printf("Key: %v", key)
	}

	value1, _ := cache.Get("1")
	log.Printf("value1: %v", value1)
	value2, _ := cache.Get("2")
	log.Printf("value2: %v", value2)

	value4, _ := cache.Get("4")
	log.Printf("value4: %v", value4)

	cache.Del("1")
	cache.Del("4")

	value1, isE := cache.Get("1")
	log.Printf("value1: %v, isE: %v", value1, isE)

	keys = cache.Keys()

	for key := range keys {
		log.Printf("Key: %v", key)
	}

}
func main() {
	Run()
}
