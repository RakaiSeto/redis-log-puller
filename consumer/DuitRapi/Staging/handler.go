package duitrapi_staging

import (
	"context"
	"database/sql"

	"github.com/rakaiseto/redis-log-puller/models"
	utils "github.com/rakaiseto/redis-log-puller/utils"
)

type DuitRapiStagingConsumer struct {
	db *sql.DB
}

func NewDuitRapiStagingConsumer(DBName string) (*DuitRapiStagingConsumer, error) {

	db, err := utils.NewDBConnectionWithCategory(models.StagingConnection, DBName)
	if err != nil {
		return nil, err
	}
	return &DuitRapiStagingConsumer{
		db: db,
	}, nil
}

func (c *DuitRapiStagingConsumer) Consume(ctx context.Context, data string) error {
	return utils.ConsumeActivityLog(ctx, c.db, data, "DuitRapi Staging")
}
