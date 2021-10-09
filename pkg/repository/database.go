package repository

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
)

// TODO

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

func (db *Database) Query(statement string) (*sql.Rows, error) {
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *Database) QueryRow(statement string) (*sql.Row, error) {
	row, err := db.QueryRow(statement)
	if err != nil {
		return nil, err
	}
	return row, err
}
