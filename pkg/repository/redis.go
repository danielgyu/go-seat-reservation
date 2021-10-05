package repository

import "github.com/go-redis/redis/v8"

type RedisDB struct {
	Client *redis.Client
}

func NewRedisClient() (*RedisDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisDB{
		Client: client,
	}, nil
}
