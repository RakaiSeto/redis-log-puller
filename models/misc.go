package models

type ConnectionCategory string

const (
	DevConnection ConnectionCategory = "DB_URL_DEV"
	StagingConnection ConnectionCategory = "DB_URL_STAGING"
	ProdConnection ConnectionCategory = "DB_URL_PROD"
)
