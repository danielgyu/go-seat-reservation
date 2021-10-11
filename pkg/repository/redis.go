package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
		log.Println("unable to connect to redis")
		panic(err)
	}

	log.Println("connected to redis")

	return &RedisDB{
		Client: client,
	}, nil
}

func (rd *RedisDB) GetItem(key string) string {
	ctx := context.Background()

	res, err := rd.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("doesn't exist in cache")
		return ""
	} else if err != nil {
		log.Println("redis GET error:", err)
		return ""
	} else {
		return res
	}
}

func (rd *RedisDB) SetItem(key string, value []byte) error {
	ctx := context.Background()

	err := rd.Client.Set(ctx, "allsHalls", value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rd *RedisDB) GetSession(token string) (int, error) {
	ctx := context.Background()

	userId, err := rd.Client.Get(ctx, fmt.Sprintf("user:%s", token)).Result()
	if err != nil {
		log.Println("no token in session:", err)
		return 0, err
	}

	id, cErr := strconv.Atoi(userId)
	if cErr != nil {
		log.Println("conversion error:", cErr)
		return 0, cErr
	}

	return id, nil
}

func (rd *RedisDB) SetSession(token string, userId int) {
	ctx := context.Background()

	if err := rd.Client.Set(ctx, fmt.Sprintf("user:%s", token), userId, 0).Err(); err != nil {
		log.Println("error settings session:", err)
	}
}
