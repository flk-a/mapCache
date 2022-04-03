package mapCache

import (
	"sync"
	"time"
)

type MapCache[cKey comparable, cValue any] struct {
	m             map[cKey]cValue
	q             []queueItem[cKey]
	ttl           time.Duration
	rear, front   int
	size, curSize int
	mu            sync.Mutex
}
type queueItem[cKey comparable] struct {
	t   int64
	val cKey
}

// NewMapCache creating new MapCache instance
func NewMapCache[cKey comparable, cValue any](size int, ttl time.Duration) *MapCache[cKey, cValue] {
	var mapCache MapCache[cKey, cValue]
	mapCache.m = make(map[cKey]cValue, size)
	mapCache.q = make([]queueItem[cKey], size)
	mapCache.ttl = ttl
	mapCache.front = -1
	mapCache.rear = -1
	mapCache.size = size

	return &mapCache
}

// Get returns value from cache. If there is no value returns found=false
func (c *MapCache[cKey, cValue]) Get(key cKey) (val cValue, found bool) {
	c.mu.Lock()
	if c.isEmpty() {
		c.mu.Unlock()
		return
	}
	c.cleanByTTL()
	v, ok := c.m[key]
	c.mu.Unlock()

	return v, ok
}

// Set writes some value to cache by key and value
func (c *MapCache[cKey, cValue]) Set(key cKey, value cValue) (val cValue) {
	c.mu.Lock()
	c.push(queueItem[cKey]{
		t:   time.Now().Add(c.ttl).UnixNano(),
		val: key,
	})
	c.m[key] = value
	c.mu.Unlock()

	return value
}

func (c *MapCache[cKey, cValue]) cleanByTTL() {
	if c.ttl == 0 {
		return
	}
	now := time.Now().UnixNano()
	for ; len(c.m) > 0; c.rear++ {
		v, found := c.top()
		if !found || now < v.t {
			break
		}
		v = c.pop()
		delete(c.m, v.val)
	}
}

func (c *MapCache[cKey, cValue]) push(qi queueItem[cKey]) {
	if c.isFull() {
		v := c.pop()
		delete(c.m, v.val)
	}
	c.curSize++
	c.rear = (c.rear + 1) % c.size
	if c.front == -1 {
		c.front = c.rear
	}
	c.q[c.rear] = qi
}

func (c *MapCache[cKey, cValue]) pop() (val queueItem[cKey]) {
	if c.isEmpty() {
		return
	}
	c.curSize--
	v := c.q[c.front]

	if c.rear == c.front {
		c.rear = -1
		c.front = -1
	} else {
		c.front = (c.front + 1) % c.size
	}

	return v
}

func (c *MapCache[cKey, cValue]) top() (val queueItem[cKey], found bool) {
	if c.isEmpty() {
		return queueItem[cKey]{}, false
	}

	return c.q[c.front], true
}

func (c *MapCache[cKey, cValue]) isFull() bool {
	return c.curSize == c.size
}

func (c *MapCache[cKey, cValue]) isEmpty() bool {
	return c.curSize == 0
}
