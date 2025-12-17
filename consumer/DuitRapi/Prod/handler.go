package duitrapi_prod

import (
	"context"
	"database/sql"

	"github.com/rakaiseto/redis-log-puller/models"
	utils "github.com/rakaiseto/redis-log-puller/utils"
)

type DuitRapiConsumer struct {
	db *sql.DB
}

func NewDuitRapiConsumer(DBName string) (*DuitRapiConsumer, error) {

	db, err := utils.NewDBConnectionWithCategory(models.ProdConnection, DBName)
	if err != nil {
		return nil, err
	}
	return &DuitRapiConsumer{
		db: db,
	}, nil
}

func (c *DuitRapiConsumer) Consume(ctx context.Context, data string) error {
	return utils.ConsumeActivityLog(ctx, c.db, data, "DuitRapi")
}
