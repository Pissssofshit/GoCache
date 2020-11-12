package lru

import (
	"fmt"
	"os"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

var lru *Cache

func setup() {
	lru = New(int64(24), nil)
}

func TestAdd(t *testing.T) {
	lru.Add("key1", String("key1"))
	lru.Add("key1", String("key1"))
	lru.Add("key2", String("key2"))
	lru.Add("key3", String("key3"))
	lru.Get("key1")
	lru.Add("key4", String("key4"))
	value, ok := lru.Get("key1")
	fmt.Printf("value of key1=\n")
	fmt.Println(value)
	if !ok {
		t.Fatalf("缓存策略有误")
	}
	_, ok = lru.Get("key2")
	if ok {
		t.Fatalf("缓存策略有误")
	}

}
func TestGet(t *testing.T) {
	value, _ := lru.Get("key1")
	fmt.Println(value)
	value, _ = lru.Get("key2")
	fmt.Println(value)
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
