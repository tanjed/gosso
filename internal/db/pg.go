package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/tanjed/go-sso/internal/config"
)

func InitPg(config *config.Config) interface{} {
	
	db, err := sql.Open("postgres", getPgConnectionString(config))

	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}
	
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	return db
}

func getPgConnectionString (config *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
}