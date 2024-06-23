package server

import (
	"fmt"

	"github.com/go-redis/redis"
)

func RedisConnection(){
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,          // use default database
	})

	err := client.Ping().Err()
	if err != nil {
		fmt.Println("Redis is not working:", err)
		return
	}
	fmt.Println("Redis is working.")
}