package utils

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"encoding/json"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rakaiseto/redis-log-puller/models"
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
	pass := url.QueryEscape(GetSecretFromKey("db", "DB_PASSWORD"))
	sqlConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", GetSecretFromKey("db", "DB_USER"), pass, GetSecretFromKey("db", "DB_HOST"), GetSecretFromKey("db", "DB_PORT"), dbName)
	fmt.Println(sqlConn)
	conn, err := sql.Open("postgres", sqlConn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return conn, nil
}

func NewDBConnectionWithCategory(category models.ConnectionCategory, dbName string) (*sql.DB, error) {
	host := GetSecretFromKey("db", "DB_HOST")
	port := GetSecretFromKey("db", "DB_PORT")
	user := GetSecretFromKey("db", "DB_USER_" + string(category))
	pass := url.QueryEscape(GetSecretFromKey("db", "DB_PASSWORD_" + string(category)))

	connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
	fmt.Println(connectionUrl)
	conn, err := sql.Open("postgres", connectionUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return conn, nil
}

func ConsumeActivityLog(
	ctx context.Context,
	db *sql.DB,
	data string,
	logPrefix string,
) error {

	var payload models.ActivityLog
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return err
	}

	_, err := db.ExecContext(
		ctx,
		`INSERT INTO activity_log
		(activity_log_id, user_id, category, activity_name,
		 entity_type, entity_id,
		 is_success, timestamp, description, metadata)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		payload.ActivityLogID,
		payload.UserID,
		payload.Category,
		payload.ActivityName,
		payload.EntityType,
		payload.EntityID,
		payload.IsSuccess,
		payload.Timestamp,
		payload.Description,
		payload.Metadata,
	)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] Processing: %s\n", logPrefix, payload.ActivityLogID)
	return nil
}