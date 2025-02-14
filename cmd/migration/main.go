package main

import (
	"log"

	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
)

func main() {
	config.Load()
	migrationQueries := db.RegisterMigrationQueries()
	database := db.Init()
	defer database.Close()
	
	

	for _, query := range migrationQueries {
	 	err := database.Conn.Query(query).Exec()

		if err != nil {
			log.Fatalf("Failed to execute query: %v\n", err)
		}
	}

}