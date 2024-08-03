package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data sync.Map
}

func NewSyncMap() Cache {
	return Cache{
		data: sync.Map{},
	}
}

func (m *Cache) Get(k string) (any, bool) {
	val, ok := m.data.Load(k)
	return val, ok
}

func (m *Cache) Set(k string, v any) {
	m.data.Store(k, v)
}

func (m *Cache) Del(k string) {
	m.data.Delete(k)
}

func (m *Cache) Contains(k string) bool {
	_, ok := m.data.Load(k)
	return ok
}

func (m *Cache) Keys() []interface{} {
	var keys []interface{}
	m.data.Range(func(key, value interface{}) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}
func main() {
	ma := NewSyncMap()
	// set key và value
	ma.Set("1", "11")
	ma.Set("2", "22")
	ma.Set("3", "33")
	ma.Set("4", "44")
	// get value of key 1
	val, ok := ma.Get("1")
	if ok {
		fmt.Println("value key 1: ", val)
	}
	// check contrains
	_, ok = ma.Get("1")
	if ok {
		fmt.Println("contrains key 1: ", val)
	}
	_, ok = ma.Get("5")
	if ok {
		fmt.Println("contrains key 5")
	} else {
		fmt.Println("not contrains key 5")
	}
	// get list keys
	list := ma.Keys()
	for _, val := range list {
		fmt.Println("key:", val)
	}
	// delete key
	ma.Del("1")
	// get value of key 1
	_, ok = ma.Get("1")
	if ok {
		fmt.Println("value key 1")
	} else {
		fmt.Println("not exis key 1")
	}
}
