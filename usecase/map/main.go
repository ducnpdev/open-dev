package main

import "log"

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
