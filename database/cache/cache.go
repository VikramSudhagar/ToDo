package cache

import (
	"github.com/go-redis/redis"
)

func CacheSetUp() *redis.Client {
	var cache = redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379",
		Password: "",
		DB:       0,
	})

	return cache
}
