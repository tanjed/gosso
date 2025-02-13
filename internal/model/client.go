package model

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/pkg/hashutilities"
)

type Client struct {
	ClientId int64 		`json:"client_id" db:"client_id"`
	ClientName string 	`json:"client_name" db:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" db:"client_secret" validate:"required"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func AuthenticateClient(clientName string, clientSecret string) bool  {
	client := getClientByClientId(clientName)

	if client == nil {
		return false
	}
	
	return hashutilities.CompareHashWithString(client.ClientSecret, clientSecret)
}

func getClientByClientId(clientName string) *Client {
	db := db.Init()
	defer db.Close()
	var client Client
	err := db.Conn.QueryRow(context.Background(), "SELECT client_name, client_secret FROM clients WHERE client_name = $1", clientName).Scan(&client.ClientName, &client.ClientSecret)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil
		}
		log.Println("Unable to fetch result", err)
		return nil
	}

	return &client
}