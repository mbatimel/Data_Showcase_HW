package cache

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/mbatimel/Data_Showcase_HW/internal/config"
)

type redisCache struct{
	redisClient	*redis.Client
	cap 		int
	mu 			sync.Mutex
	keys  		map[interface{}]time.Time
}


func NewRedisCache(cfg config.Cache) (ICache, error){
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		DB:       0,
	})

	err := client.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &redisCache{
		redisClient:	client,
		cap:			cfg.Cap,
		mu:				sync.Mutex{},
		keys:			make(map[interface{}]time.Time),
	}, nil
}

func (cs *redisCache) Cap() int {
	return cs.cap
}

func (cs *redisCache) Clear() error{
	cs.mu.Lock()
	defer cs.mu.Unlock()
	err :=cs.redisClient.FlushDB().Err()
	if err != nil {
		return err
	}
	clear(cs.keys)
	return nil
}

func (cs *redisCache) Get(key interface{}) (value interface{}, ok bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    val:= cs.redisClient.Get(fmt.Sprintf("%v", key)).Val()
    if val == "" {
        return nil, false
    } 
	cs.keyUsedUpdate(key)
    return val, true
}

func (cs *redisCache) Remove(key interface{}) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    err := cs.redisClient.Del(fmt.Sprintf("%v", key)).Err()
	if err != nil{
		log.Fatalln(err)
	}
	delete(cs.keys, key)
}

func (cs *redisCache) Add(key, value interface{}) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    cs.addItem(key,value,0)
}

func (cs *redisCache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
    cs.addItem(key,value,ttl)
}

func (cs *redisCache) keyUsedUpdate(key interface{}){
	cs.keys[key] = time.Now()
}

func (cs *redisCache) addItem(key, value interface{}, ttl time.Duration) {
	if len(cs.keys) >= cs.cap {
		cs.evictLRU()
	}
	cs.redisClient.Set(fmt.Sprintf("%v", key), value, ttl)
	cs.keyUsedUpdate(key)
}

func (cs *redisCache) evictLRU() {
	// Найти ключ с наименьшим временем последнего доступа
	var lruKey interface{}
	var lruTime time.Time
	first := true

	for key, accessTime := range cs.keys {
		if first || accessTime.Before(lruTime) {
			lruKey = key
			lruTime = accessTime
			first = false
		}
	}

	// Удалить найденный LRU ключ из Redis и из отслеживаемых ключей
	if lruKey != nil {
		err := cs.redisClient.Del(fmt.Sprintf("%v", lruKey)).Err()
		if err != nil {
			log.Printf("Failed to delete LRU key %v: %v", lruKey, err)
		}
		delete(cs.keys, lruKey)
	}
}
