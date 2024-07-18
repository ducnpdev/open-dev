package main

import (
	"log"
	"sync"
)

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

func Run() {
	cache := NewCache()

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
