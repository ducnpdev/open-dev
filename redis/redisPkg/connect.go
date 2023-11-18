package redisPkg

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// var RedisIn *redis.Client

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return rdb
}
