package duitrapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	utils "github.com/rakaiseto/redis-log-puller/utils"
)

type DuitRapiConsumer struct {
	db *sql.DB
}

func NewDuitRapiConsumer(DBName string) (*DuitRapiConsumer, error) {

	db, err := utils.NewDBConnection(DBName)
	if err != nil {
		return nil, err
	}
	return &DuitRapiConsumer{
		db: db,
	}, nil
}

func (c *DuitRapiConsumer) Consume(ctx context.Context, data string) error {
	// unmarshal json
	var ocrMarketplaceData map[string]any
	if err := json.Unmarshal([]byte(data), &ocrMarketplaceData); err != nil {
		return err
	}

	// insert into db
	_, err := c.db.ExecContext(
		ctx,
		`INSERT INTO activity_log (activity_log_id, user_id, category, activity_name, is_success, timestamp, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ocrMarketplaceData["activity_log_id"],
		ocrMarketplaceData["user_id"],
		ocrMarketplaceData["category"],
		ocrMarketplaceData["activity_name"],
		ocrMarketplaceData["is_success"],
		ocrMarketplaceData["timestamp"],
		ocrMarketplaceData["description"],
	)
	if err != nil {
		return fmt.Errorf("failed to insert DuitRapi data: %w", err)
	}

	fmt.Printf("[DuitRapi] Processing: %s\n", ocrMarketplaceData["activity_log_id"])

	return nil
}
