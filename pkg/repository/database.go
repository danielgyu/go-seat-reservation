package repository

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

type Database struct {
	rdbms *sql.DB
	cache *redis.Client
}

func (db *Database) SetString(ctx context.Context, key, value string) {
	db.cache.Set(ctx, key, value, 0)
}

func (db *Database) GetString(ctx context.Context, key string) (string, error) {
	val, err := db.cache.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
