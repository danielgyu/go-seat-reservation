package repository

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

func (rd *RedisDB) SetString(ctx context.Context, key string, value string) {
	rd.Client.Set(ctx, key, value, 0)
}

func (rd *RedisDB) GetString(ctx context.Context, key string) (string, error) {
	val, err := rd.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
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
