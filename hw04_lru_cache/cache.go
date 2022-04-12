package hw04lrucache

import "sync"

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
	m        sync.Mutex
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

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.m.Lock()
	listItem, ok := l.items[key]
	if ok {
		cacheItem := listItem.Value.(*cacheItem)
		cacheItem.value = value
		l.queue.MoveToFront(listItem)
		l.m.Unlock()
		return true
	}
	if l.queue.Len() == l.capacity {
		lastItem := l.queue.Back()
		lastKey := lastItem.Value.(*cacheItem).key
		delete(l.items, lastKey)
		l.queue.Remove(lastItem)
	}
	newListItem := l.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	l.items[key] = newListItem
	l.m.Unlock()
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.m.Lock()
	listItem, ok := l.items[key]
	if ok {
		cacheItem := listItem.Value.(*cacheItem)
		l.queue.MoveToFront(listItem)
		l.m.Unlock()
		return cacheItem.value, true
	}
	l.m.Unlock()
	return nil, false
}

func (l *lruCache) Clear() {
	l.m.Lock()
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
	l.m.Unlock()
}
