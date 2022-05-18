package cache

import "github.com/gofiber/storage/redis"

func CacheSetUp() *redis.Storage {
	store := redis.New(redis.Config{
		Host:      "host.docker.internal",
		Port:      6379,
		Username:  "",
		Password:  "",
		URL:       "",
		Database:  0,
		Reset:     true,
		TLSConfig: nil,
	})

	return store
}
