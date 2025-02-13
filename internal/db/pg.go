package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/tanjed/go-sso/internal/config"
)

func initPg(config *config.Config) *pgx.Conn {

	db, err := pgx.Connect(context.Background(), getPgConnectionString(config))

	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}
	
	return db
}

func closePg(db *pgx.Conn) {
	db.Close(context.Background())
}

func getPgConnectionString (config *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
}