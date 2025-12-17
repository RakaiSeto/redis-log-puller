package duitrapi_dev

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/rakaiseto/redis-log-puller/models"
	utils "github.com/rakaiseto/redis-log-puller/utils"
)

type DuitRapiDevConsumer struct {
	db *sql.DB
}

func NewDuitRapiDevConsumer(DBName string) (*DuitRapiDevConsumer, error) {

	db, err := utils.NewDBConnectionWithCategory(models.DevConnection, DBName)
	if err != nil {
		return nil, err
	}
	return &DuitRapiDevConsumer{
		db: db,
	}, nil
}

func (c *DuitRapiDevConsumer) Consume(ctx context.Context, data string) error {
	// unmarshal json
	var duitrapiData models.ActivityLog
	if err := json.Unmarshal([]byte(data), &duitrapiData); err != nil {
		return err
	}

	// insert into db
	_, err := c.db.ExecContext(
		ctx,
		`INSERT INTO activity_log (activity_log_id, user_id, category, activity_name, is_success, timestamp, description, metadata) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		duitrapiData.ActivityLogID,
		duitrapiData.UserID,
		duitrapiData.Category,
		duitrapiData.ActivityName,
		duitrapiData.IsSuccess,
		duitrapiData.Timestamp,
		duitrapiData.Description,
		duitrapiData.Metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to insert DuitRapi Dev data: %w", err)
	}

	fmt.Printf("[DuitRapi Dev] Processing: %s\n", duitrapiData.ActivityLogID)

	return nil
}
