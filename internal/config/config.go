package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port int `json:"PORT"`
	DB_HOST string `json:"DB_HOST"`
	DB_PORT int `json:"DB_PORT"`
	DB_USER string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_NAME string `json:"DB_NAME"`
	DB_DRIVER string `json:"DB_DRIVER"`
	JWT_SECRET string `json:"JWT_SECRET"`
	REDIS_HOST string `json:"REDIS_HOST"`
	REDIS_PORT int `json:"REDIS_PORT"`
}


func NewConfig() *Config{
	var appConfig Config
	content, err := os.ReadFile("internal/config/config.json")

	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(content, &appConfig)

	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return &appConfig
}

func (c *Config) Close() {
	c = &Config{}
}