package lrucache

import (
	"errors"
	"sync"
	"time"
)

// Lrucache is a cache structure
type Lrucache struct {
	sync.Mutex
	capacity   int
	cache      map[string]*entry
	firstEntry *entry
	lastEntry  *entry
}

type entry struct {
	key   string
	value Value
	added time.Time

	nextEntry *entry
	prevEntry *entry
}

// Value is a value type
type Value interface{}

// Cache is a interface of cache
type Cache interface {
	Get(key string) (Value, bool)
	Set(key string, value Value)
	Del(key string)
	Len() int
	Flush()
}

// NewLrucache returns new Lrucache cache structure
func NewLrucache(capacity int) (Cache, error) {
	if capacity <= 0 {
		return nil, errors.New("invalid capacity")
	}
	return &Lrucache{
		capacity:   capacity,
		cache:      make(map[string]*entry),
		firstEntry: nil,
		lastEntry:  nil,
	}, nil
}

func (lru *Lrucache) set(key string, value Value) {
	e := &entry{
		key:       key,
		value:     value,
		added:     time.Now(),
		prevEntry: nil,
		nextEntry: nil,
	}

	if lru.Len() == lru.capacity && lru.lastEntry != nil {
		nextLastEntry := lru.lastEntry.nextEntry
		lru.del(lru.lastEntry.key)
		lru.lastEntry = nextLastEntry
	}
	if lru.firstEntry == nil && lru.lastEntry == nil {
		lru.cache[key] = e
		lru.firstEntry = e
		lru.lastEntry = e
	}
	if lru.firstEntry != nil {
		e.prevEntry = lru.firstEntry
		lru.firstEntry.nextEntry = e
		lru.firstEntry = e
		lru.cache[key] = e
	}
}

// Set creates new entry in the cache
func (lru *Lrucache) Set(key string, value Value) {
	lru.Lock()
	defer lru.Unlock()

	lru.set(key, value)
}

// Get returns value from the cache
func (lru *Lrucache) Get(key string) (Value, bool) {
	lru.Lock()
	defer lru.Unlock()

	entry, ok := lru.cache[key]
	if !ok {
		return nil, false
	}
	return entry.value, true
}

func (lru *Lrucache) del(key string) {
	e := lru.cache[key]
	if e == nil {
		return
	}

	if e.nextEntry != nil {
		e.nextEntry.prevEntry = e.prevEntry
	}
	if e.prevEntry != nil {
		e.prevEntry.nextEntry = e.nextEntry
	}
	delete(lru.cache, key)
}

// Del removes entry from the cache
func (lru *Lrucache) Del(key string) {
	lru.Lock()
	defer lru.Unlock()

	lru.del(key)
}

// Len returns length of the cache
func (lru *Lrucache) Len() int {
	return len(lru.cache)
}

// Flush cleans the cache
func (lru *Lrucache) Flush() {
	for k := range lru.cache {
		delete(lru.cache, k)
	}
}
