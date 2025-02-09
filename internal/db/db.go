package db

import (
	"github.com/tanjed/go-sso/internal/config"
)

const DB_DRIVER_POSTGRES = "postgres"

var DB interface{}

func Init() {
	config := config.AppConfig

	switch config.DB_DRIVER {
	case DB_DRIVER_POSTGRES :
		DB = InitPg(&config)
	default:
		DB = InitPg(&config)
	}
}