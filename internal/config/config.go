package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port string `json:"PORT"`
	DB_HOST string `json:"DB_HOST"`
	DB_PORT string `json:"DB_PORT"`
	DB_USER string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_NAME string `json:"DB_NAME"`
	DB_DRIVER string `json:"DB_DRIVER"`
}

var AppConfig Config

func Load() {
		
	content, err := os.ReadFile("internal/config/config.json")

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(content, &AppConfig)

	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
}