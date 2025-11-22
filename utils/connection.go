package utils

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", GetSecretFromKey("db", "REDIS_HOST"), GetSecretFromKey("db", "REDIS_PORT")),
		Password: GetSecretFromKey("db", "REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}

func NewDBConnection(dbName string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", GetSecretFromKey("db", "DB_USER"), GetSecretFromKey("db", "DB_PASSWORD"), GetSecretFromKey("db", "DB_HOST"), GetSecretFromKey("db", "DB_PORT"), dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return conn, nil
}
