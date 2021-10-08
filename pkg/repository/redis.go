package repository

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

func NewRedisClient() (*RedisDB, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	return &RedisDB{
		Client: client,
	}, nil
}
