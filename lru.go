/*
	- Set: Add the item both in queue and HashMap. If they capacity is full, it removes the least recently used
	element

	- Get: Returns the item requested via Key. On querying the item it comes to forward of the queue

*/

package main

import (
	"container/list"
	"errors"
	"fmt"
)

var queue = list.New()
var m = make(map[string]string)

/*
	Cache struct
*/
type Cache struct {
	capacity int
}

/*
	CacheItem Struct
*/
type CacheItem struct {
	Name  string
	Value string
}

// Move the list item to front
func moveListElement(k string) {
	for e := queue.Front(); e != nil; e = e.Next() {
		if e.Value.(CacheItem).Name == k {
			queue.MoveToFront(e)
			break
		}
	}
}

func (cache *Cache) print() {

	fmt.Println("Printing Queue Items")
	// Iterate through list and print its contents.
	for e := queue.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(CacheItem).Name)
	}
}

func (cache *Cache) set(key string, val string) string {
	//Search the key in map
	_, found := m[key]

	if !found {
		//Check the capacity
		if len(m) == cache.capacity { // Time to evict
			// Get the least use item from the queue
			e := queue.Back()
			queue.Remove(e) // Dequeue
			keyName := e.Value.(CacheItem).Name
			// Delete from the map
			delete(m, keyName)
		} else {
			//There is still some room
			item := CacheItem{Name: key, Value: val}
			queue.PushFront(item)
			m[key] = val
		}
	}
	return "1"
}

func (cache *Cache) get(k string) (string, error) {
	//Search the key in map
	v, found := m[k]
	if found {
		v := m[k]
		//fmt.Println(v)
		moveListElement(v)
		return v, nil
	}
	v = "-1"
	return v, errors.New("Key not found")
}
