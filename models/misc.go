package models

type ConnectionCategory string

const (
	DevConnection ConnectionCategory = "DEV"
	StagingConnection ConnectionCategory = "STAGING"
	ProdConnection ConnectionCategory = "PROD"
)
