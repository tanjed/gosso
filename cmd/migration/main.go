package main

import (
	"fmt"

	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/model"
)

func main() {
	fmt.Println("Initiating client seeder")
	config.NewConfig()
	client := model.NewClient("search_service", "secret")
	fmt.Println(client)
	fmt.Println(client.Insert())
	
	// migrationQueries := db.RegisterMigrationQueries()
	// database := db.InitDB()

	// for _, query := range migrationQueries {
	//  	err := database.Conn.Query(query).Exec()

	// 	if err != nil {
	// 		log.Fatalf("Failed to execute query: %v\n", err)
	// 	}
	// }

}