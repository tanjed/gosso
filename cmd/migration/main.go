package main

import (
	"context"
	"log"

	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
)

func main() {
	config.Load()
	migrationQueries := db.RegisterMigrationQueries()
	database := db.Init()
	defer database.Close()
	tx, err := database.Conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Unable to begin transaction: %v\n", err)
	}

	defer tx.Rollback(context.Background())

	for _, query := range migrationQueries {
		_, err := database.Conn.Exec(context.Background(), query)

		if err != nil {
			log.Fatalf("Failed to execute CREATE query: %v\n", err)
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		log.Fatalf("Failed to commit transaction: %v\n", err)
	}

}