package service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis"
)
type ICache interface {
	Cap() int
	Clear()
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

type CacheService struct{
	redisClient	*redis.Client
	cap 		int
	mu 			sync.Mutex
	keys  		map[interface{}]time.Time
}


func NewRedisConnection( cap int) (*CacheService, error){
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
	})

	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &CacheService{
		redisClient:	client,
		cap:			cap,
		mu:				sync.Mutex{},
		keys:			make(map[interface{}]time.Time),
	}, nil
}

func (cs *CacheService) Cap() int {
	return cs.cap
}
func (cs *CacheService) Clear() error{
	cs.mu.Lock()
	defer cs.mu.Unlock()
	err :=cs.redisClient.FlushDB().Err()
	if err != nil {
		return err
	}
	clear(cs.keys)
	return nil
}
func (cs *CacheService) Get(key interface{}) (value interface{}, ok bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    val:= cs.redisClient.Get(fmt.Sprintf("%v", key)).Val()
    if val == "" {
        return nil, false
    } 
	cs.keyUsedUpdate(key)
    return val, true
}

func (cs *CacheService) Remove(key interface{}) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    err := cs.redisClient.Del(fmt.Sprintf("%v", key)).Err()
	if err != nil{
		log.Fatalln(err)
	}
	delete(cs.keys, key)
}

func (cs *CacheService) Add(key, value interface{}) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    cs.addItem(key,value,0)
}
func (cs *CacheService) AddWithTTL(key, value interface{}, ttl time.Duration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    cs.addItem(key,value,ttl)
}

func (cs *CacheService) keyUsedUpdate(key interface{}){
	cs.keys[key] = time.Now()
}

func (cs *CacheService) addItem(key, value interface{}, ttl time.Duration) {
	if len(cs.keys) >= cs.cap {
		cs.evictLRU()
	}
	cs.redisClient.Set(fmt.Sprintf("%v", key), value, ttl * time.Second)
	cs.keyUsedUpdate(key)
}

func (cs *CacheService) evictLRU() {
	var oldestKey interface{}
	var oldestTime time.Time

	for k, v := range cs.keys {
		if oldestKey == nil || v.Before(oldestTime) {
			oldestKey = k
			oldestTime = v
		}
	}
	cs.Remove(oldestKey)
}