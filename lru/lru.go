package lru

import (
	"container/list"
	"fmt"
)

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64                    //最大存储大小
	nbytes   int64                    //当前存储大小
	ll       *list.List               //维护的队列
	cache    map[string]*list.Element //字典数据
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// 每一个队列节点存储的基本元素
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 增、改
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		entryItem := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(entryItem.value.Len())
		entryItem.value = value
	} else {
		entryItem := c.ll.PushFront(&entry{key, value})
		c.cache[key] = entryItem
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	if c.maxBytes < c.nbytes {
		fmt.Println("清除老数据")
		c.RemoveOldest()
	}
	fmt.Printf("当前缓存数据容量%d,已存储%d\n", c.maxBytes, c.nbytes)
}

// 删
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		entryEle := ele.Value.(*entry)
		delete(c.cache, entryEle.key)
		fmt.Printf("当前已占内存%d\n", c.nbytes)
		c.nbytes = c.nbytes - int64(entryEle.value.Len()) - int64(len(entryEle.key))
		fmt.Printf("本次清除内存%d\n", int64(entryEle.value.Len()))
	}
}

// 查
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		entryEle := ele.Value.(*entry)
		return entryEle.value, ok
	}
	return
}
