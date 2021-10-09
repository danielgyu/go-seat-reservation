package repository

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
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

func CheckCache(h httprouter.Handle, rd *RedisDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		hallId, err := strconv.Atoi(param.ByName("id"))
		if err != nil {
			log.Println("path param error:", err)
		}

		ctx := context.Background()
		res, err := rd.Client.Get(ctx, fmt.Sprintf("hall:%d", hallId)).Result()
		if err == nil {
			fmt.Fprintf(w, res)
			return
		}

		h(w, r, param)
	}
}
