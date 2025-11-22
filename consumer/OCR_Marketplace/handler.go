package ocr_marketplace

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	utils "github.com/rakaiseto/redis-log-puller/utils"
)

type OCRMarketplaceConsumer struct {
	db *sql.DB
}

func NewOCRMarketplaceConsumer(DBName string) (*OCRMarketplaceConsumer, error) {

	db, err := utils.NewDBConnection(DBName)
	if err != nil {
		return nil, err
	}
	return &OCRMarketplaceConsumer{
		db: db,
	}, nil
}

func (c *OCRMarketplaceConsumer) Consume(ctx context.Context, data string) error {
	// unmarshal json
	var ocrMarketplaceData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &ocrMarketplaceData); err != nil {
		return err
	}

	// insert into db
	_, err := c.db.ExecContext(
		ctx,
		`INSERT INTO activity_log (activity_log_id, account_id, category, activity_name, is_success, timestamp, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ocrMarketplaceData["activity_log_id"],
		ocrMarketplaceData["account_id"],
		ocrMarketplaceData["category"],
		ocrMarketplaceData["activity_name"],
		ocrMarketplaceData["is_success"],
		ocrMarketplaceData["timestamp"],
		ocrMarketplaceData["description"],
	)
	if err != nil {
		return fmt.Errorf("failed to insert OCR marketplace data: %w", err)
	}

	fmt.Printf("[OCRMarketplace] Processing: %s\n", ocrMarketplaceData["activity_log_id"])

	return nil
}
