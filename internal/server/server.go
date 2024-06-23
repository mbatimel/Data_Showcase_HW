package server

import (
	"context"
	"fmt"
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
	redisClient *redis.Client
	cap int
	ctx context.Context
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
		redisClient: client,
		cap: cap,
		ctx: context.Background(),
	}, nil
}

func (cs *CacheService) Cap() int {
	return cs.cap
}
func (cs *CacheService) Clear() error{
	err :=cs.redisClient.Close()
	if err != nil {
		return err
	}
	return nil
}
func (cs *CacheService) Get(key interface{}) (value interface{}, ok bool) {
    val:= cs.redisClient.Get(fmt.Sprintf("%v", key)).Val()
    if val == "" {
        return nil, false
    } 
    return val, true
}

func (cs *CacheService) Remove(key interface{}) {
    cs.redisClient.Del(fmt.Sprintf("%v", key))
}
func (cs *CacheService) Add(key, value interface{}) {
    cs.redisClient.Set(fmt.Sprintf("%v", key), value, 0)
}
func (cs *CacheService) AddWithTTL(key, value interface{}, ttl time.Duration) {
    cs.redisClient.Set(fmt.Sprintf("%v", key), value, ttl)
}