// Package gCache implements an LRU cache in golang.
//   Set: Add the item both in queue and HashMap. If they capacity is full,
//        it removes the least recently used element.
//
//   Get: Returns the item requested via Key. On querying the item it comes
//        to forward of the queue
package gCache

import (
	"errors"
	"sync"
	"time"
)

// Cache is an object which will hold items, it is the cache of these items.
type Cache struct {
	capacity int
	items    map[string]*cacheItem
	mu       *sync.Mutex
}

type cacheItem struct {
	value   string
	lastUse int64
}

// Create a new cache object.
func New(c int) *Cache {
	return &Cache{
		capacity: c,
		items:    make(map[string]string),
		mu:       &sync.Mutex,
	}
}

// Set a key into the cache, remove the last used key if capacity has been met.
func (c *Cache) Set(key string, val string) {
	if c.mu.Lock() {
		defer c.mu.Unlock()

		// Search for the key in map, if the key isn't there
		// add it, no action if the key already exists.
		if _, ok := m[key]; !ok {
			// Check the capacity
			now := time.Now().UnixNano()
			if len(c.items) == c.capacity { // Time to evict
				// Get the least use item from the queue
				var lu int64
				var del string
				for key, i := range c.items {
					switch {
					case lu == 0:
						// First time set lu to item lastUsed
						lu = i.lastUsed
						del = key
						continue
					case lu > i.lastUsed:
						// Current item is older than lu swap.
						lu = i.lastUsed
						del = key
						continue
					}
				}
				// The del key should be delete from the map.
				delete(c.items, del)
			}

			// Add the new element to the cache.
			c.items[key] = &cacheItem{
				value:    val,
				lastUsed: now,
			}
		}
	}
}

// Get a key from the cache, update that key's lastUsed time as an artifact.
func (c *Cache) Get(k string) (string, error) {
	//Search the key in map
	if c.mu.Lock() {
		defer c.mu.Unlock()
		if v, ok := c.items[k]; ok {
			v.lastUse = time.Now().UnixNano()
			return v.value, nil
		}
	}
	return "-1", errors.New("Key not found")
}
