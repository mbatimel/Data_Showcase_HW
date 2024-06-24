package cache

import (
	"container/list"
	"sync"
	"time"
)

type inMemoryCache struct {
	capacity int
	cache    map[interface{}]*list.Element
	ll       *list.List
	mu       sync.Mutex
}

type entry struct {
	key   interface{}
	value interface{}
	ttl   *time.Time
}

func NewInMemoryCache(capacity int) ICache {
	return &inMemoryCache{
		capacity: capacity,
		cache:    make(map[interface{}]*list.Element),
		ll:       list.New(),
	}
}

func (imc *inMemoryCache) Cap() int {
	return imc.capacity
}

func (imc *inMemoryCache) Clear() error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	imc.cache = make(map[interface{}]*list.Element)
	imc.ll.Init()
	return nil
}

func (imc *inMemoryCache) Get(key interface{}) (interface{}, bool) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		if imc.isExpired(elem.Value.(*entry)) {
			imc.removeElement(elem)
			return nil, false
		}
		imc.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (imc *inMemoryCache) Remove(key interface{}) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.removeElement(elem)
	}
}

func (imc *inMemoryCache) Add(key, value interface{}) {
	imc.addItem(key, value, nil)
}

func (imc *inMemoryCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl)
	imc.addItem(key, value, &expiration)
}

func (imc *inMemoryCache) addItem(key, value interface{}, ttl *time.Time) {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	if elem, ok := imc.cache[key]; ok {
		imc.ll.MoveToFront(elem)
		elem.Value.(*entry).value = value
		elem.Value.(*entry).ttl = ttl
		return
	}

	if imc.ll.Len() >= imc.capacity {
		imc.evictLRU()
	}

	ent := &entry{key, value, ttl}
	elem := imc.ll.PushFront(ent)
	imc.cache[key] = elem
}

func (imc *inMemoryCache) evictLRU() {
	elem := imc.ll.Back()
	if elem != nil {
		imc.removeElement(elem)
	}
}

func (imc *inMemoryCache) removeElement(elem *list.Element) {
	imc.ll.Remove(elem)
	delete(imc.cache, elem.Value.(*entry).key)
}

func (imc *inMemoryCache) isExpired(ent *entry) bool {
	if ent.ttl == nil {
		return false
	}
	return ent.ttl.Before(time.Now())
}
