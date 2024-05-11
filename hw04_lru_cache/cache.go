package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	qItem, ok := cache.items[key]
	if !ok {
		return nil, false
	}

	cache.queue.MoveToFront(qItem)

	return getCacheItem(qItem).value, true
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	qItem, ok := cache.items[key]
	if ok {
		qItem.Value = cacheItem{
			key:   key,
			value: value,
		}
		cache.queue.MoveToFront(qItem)
		return true
	}

	cItem := cacheItem{
		key:   key,
		value: value,
	}
	qItem = cache.queue.PushFront(cItem)
	cache.items[key] = qItem

	if cache.queue.Len() > cache.capacity {
		qLast := cache.queue.Back()
		cLast := getCacheItem(qLast)
		delete(cache.items, cLast.key)
		cache.queue.Remove(qLast)
	}

	return false
}

func (cache *lruCache) Clear() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}

func getCacheItem(queueItem *ListItem) cacheItem {
	cacheItem, ok := queueItem.Value.(cacheItem)
	if !ok {
		panic("unexpected value type")
	}

	return cacheItem
}
