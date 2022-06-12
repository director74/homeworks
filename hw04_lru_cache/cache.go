package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	liP, ok := lru.items[key]
	if ok {
		liP.Value = cacheItem{
			key:   key,
			value: value,
		}
		lru.queue.MoveToFront(liP)

		return true
	}

	if lru.queue.Len() >= lru.capacity {
		last := lru.queue.Back()
		delete(lru.items, last.Value.(cacheItem).key)
		lru.queue.Remove(last)
	}
	lru.items[key] = lru.queue.PushFront(cacheItem{
		key:   key,
		value: value,
	})

	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	liP, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(liP)
		return liP.Value.(cacheItem).value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	for lru.queue.Back() != nil {
		delete(lru.items, lru.queue.Back().Value.(cacheItem).key)
		lru.queue.Remove(lru.queue.Back())
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
