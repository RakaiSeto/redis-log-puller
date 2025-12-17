package duitrapi_dev

import (
	"context"
	"database/sql"

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
	return utils.ConsumeActivityLog(ctx, c.db, data, "DuitRapi Dev")
}
