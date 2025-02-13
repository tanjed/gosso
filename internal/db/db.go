package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/tanjed/go-sso/internal/config"
)

const DB_DRIVER_POSTGRES = "postgres"

type DB struct {
	Conn *pgx.Conn
}

func Init() DB{
	var db DB
	switch config.AppConfig.DB_DRIVER {
	case DB_DRIVER_POSTGRES :
		db.Conn = initPg(&config.AppConfig)
	default:
		db.Conn = initPg(&config.AppConfig)
	}

	return db
}

func (db *DB) Close() {
	switch config.AppConfig.DB_DRIVER {
	case DB_DRIVER_POSTGRES :
		closePg(db.Conn)
		
	default:
		closePg(db.Conn)
	}
}