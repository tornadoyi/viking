package cache

import (
	"fmt"
	"sync"
)

// Cache is a thread-safe fixed size LRU cache.
type CLRU struct {
	lru  *LRU
	lock sync.RWMutex
}

// New creates an LRU of the given size.
func New(size int) *CLRU {
	return &CLRU {
		NewLRU(size),
		sync.RWMutex{},
	}
}

// Len returns the number of items in the cache.
func (c *CLRU) Len() int {
	c.lock.RLock()
	length := c.lru.Len()
	c.lock.RUnlock()
	return length
}

// Contains checks if a key is in the cache, without updating the
// recent-ness or deleting it for being stale.
func (c *CLRU) Contains(key interface{}) bool {
	c.lock.RLock()
	containKey := c.lru.Contains(key)
	c.lock.RUnlock()
	return containKey
}

// ContainsOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *CLRU) ContainsOrAdd(key, value interface{}) (ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lru.Contains(key) {
		return true, false
	}
	evicted = c.lru.Add(key, value)
	return false, evicted
}

// Get looks up a key's value from the cache.
func (c *CLRU) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	value, ok = c.lru.Get(key)
	c.lock.Unlock()
	return value, ok
}

// Get looks up a key's value from the cache.
func (c *CLRU) Gets(keys []interface{}) ([]interface{}, []bool) {
	c.lock.Lock()
	values, oks := make([]interface{}, len(keys)), make([]bool, len(keys))
	for i, key := range keys {
		value, ok := c.lru.Get(key)
		values[i], oks[i] = value, ok
	}
	c.lock.Unlock()
	return values, oks
}


// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *CLRU) Peek(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	value, ok = c.lru.Peek(key)
	c.lock.RUnlock()
	return value, ok
}


// Peek returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *CLRU) Peeks(keys []interface{}) ([]interface{}, []bool) {
	c.lock.Lock()
	values, oks := make([]interface{}, len(keys)), make([]bool, len(keys))
	for i, key := range keys {
		value, ok := c.lru.Peek(key)
		values[i], oks[i] = value, ok
	}
	c.lock.Unlock()
	return values, oks
}


// PeekLatest returns the latest entry
func (c *CLRU) PeekLatest() (key interface{}, value interface{}, ok bool) {
	c.lock.RLock()
	key, value, ok = c.lru.PeekLatest()
	c.lock.RUnlock()
	return key, value, ok
}

// PeekOldest returns the oldest entry
func (c *CLRU) PeekOldest() (key interface{}, value interface{}, ok bool) {
	c.lock.RLock()
	key, value, ok = c.lru.PeekOldest()
	c.lock.RUnlock()
	return key, value, ok
}


// PeekOrAdd checks if a key is in the cache without updating the
// recent-ness or deleting it for being stale, and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (c *CLRU) PeekOrAdd(key, value interface{}) (previous interface{}, ok, evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	previous, ok = c.lru.Peek(key)
	if ok {
		return previous, true, false
	}

	evicted = c.lru.Add(key, value)
	return nil, false, evicted
}

// Purge is used to completely clear the cache.
func (c *CLRU) Purge() {
	c.lock.Lock()
	c.lru.Purge()
	c.lock.Unlock()
}

// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *CLRU) Add(key, value interface{}) (evicted bool) {
	c.lock.Lock()
	evicted = c.lru.Add(key, value)
	c.lock.Unlock()
	return evicted
}


// Add adds a value to the cache. Returns true if an eviction occurred.
func (c *CLRU) Adds(keys []interface{}, values[]interface{}) ([]bool, error) {
	if len(keys) != len(values) { return nil, fmt.Errorf("mismatch length keys: %v values: %v", len(keys), len(values))}
	c.lock.Lock()
	evicteds := make([]bool, len(keys))
	for i:=0; i<len(keys); i++ {
		key, value := keys[i], values[i]
		evicted := c.lru.Add(key, value)
		evicteds[i] = evicted
	}
	c.lock.Unlock()
	return evicteds, nil
}


// Remove removes the provided key from the cache.
func (c *CLRU) Remove(key interface{}) (present bool) {
	c.lock.Lock()
	present = c.lru.Remove(key)
	c.lock.Unlock()
	return
}


// Remove removes the provided key from the cache.
func (c *CLRU) Removes(keys []interface{}) ([]bool) {
	c.lock.Lock()
	presents := make([]bool, len(keys))
	for i, key := range keys {
		present := c.lru.Remove(key)
		presents[i] = present
	}
	c.lock.Unlock()
	return presents
}


// Resize changes the cache size.
func (c *CLRU) Resize(size int) (evicted int) {
	c.lock.Lock()
	evicted = c.lru.Resize(size)
	c.lock.Unlock()
	return evicted
}


// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *CLRU) Keys() []interface{} {
	c.lock.RLock()
	keys := c.lru.Keys()
	c.lock.RUnlock()
	return keys
}


