package model

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/pkg/hashutilities"
)

const TIME_FORMAT = "2006-01-02 15:04:05.0000"
type Client struct {
	ClientId string 	`json:"client_id" db:"client_id"`
	ClientName string 	`json:"client_name" db:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" db:"client_secret" validate:"required"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type clientUnAuthorizationCode struct {
	ClientId string
	ClientCode string
	RedirectUri string
	GeneratedAt time.Time
	ExpiredAt time.Time	
}

type ClientNotFoundError struct {
	Message string
	Code int
}

type ClientUnAuthorizedError struct {
	Message string
	Code int
}

func (e ClientNotFoundError) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

func (e ClientUnAuthorizedError) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

func AuthenticateClient(clientName string, clientSecret string) (*Client, error)  {
	c := getClientByClientId(clientName)
	
	if c == nil {
		return nil, &ClientNotFoundError{
			Message: "Client not found",
			Code: http.StatusNotFound,
		}
	}

	client := *c
	if !hashutilities.CompareHashWithString(client.ClientSecret, clientSecret) {
		return nil, &ClientNotFoundError{
			Message: "Client unauthorized",
			Code: http.StatusUnauthorized,
		}
	}

	return c, nil 
}

func getClientByClientId(clientName string) *Client {
	db := db.Init()
	defer db.Close()
	var client Client
	err := db.Conn.Query("SELECT client_id, client_name, client_secret, created_at, updated_at  FROM clients WHERE client_name = ?", clientName).
	Scan(&client.ClientId, &client.ClientName, &client.ClientSecret, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil
	}

	return &client
}

func ClientHasValidSession(clientId string, redirectUri string) (*clientUnAuthorizationCode, error) {
	db := db.Init()
	defer db.Close()
	var clientAuthorizationCode clientUnAuthorizationCode
	err := db.Conn.Query("SELECT client_id, client_code, redirect_uri, generated_at, expired_at FROM client_authorization_codes WHERE client_id = ? AND redirect_uri = ? AND expired_at > ?", clientId, redirectUri, time.Now()).
	Scan(&clientAuthorizationCode.ClientId, &clientAuthorizationCode.ClientCode, &clientAuthorizationCode.RedirectUri, &clientAuthorizationCode.GeneratedAt, &clientAuthorizationCode.ExpiredAt)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil, err
	}

	return &clientAuthorizationCode, nil
}

func NewClientAuthorizationCode(clientId string, redirectUri string) (*clientUnAuthorizationCode, error){
	clientUnAuthorizationCode := clientUnAuthorizationCode {
		ClientId: clientId,
		ClientCode: generateAuthorizationCode(),
		RedirectUri: redirectUri,
		GeneratedAt: time.Now(),
		ExpiredAt: time.Now().Add(30 * time.Second),
	}
	err := insertClientAuthorizationCode(clientUnAuthorizationCode)

	if err != nil {
		return nil, err
	}
	return &clientUnAuthorizationCode, nil
}

func insertClientAuthorizationCode(c clientUnAuthorizationCode) error {
	db := db.Init()
	defer db.Close()

	return db.Conn.Query("INSERT INTO client_authorization_codes (client_id, client_code, redirect_uri, generated_at, expired_at) VALUES (?,?,?,?,?)", c.ClientId, c.ClientCode, c.RedirectUri, c.GeneratedAt, c.ExpiredAt).Exec()
}

func generateAuthorizationCode() string {
	code := make([]byte, 32)
	_, err := rand.Read(code)
	if err != nil {
		log.Fatal("Failed to generate authorization code:", err)
	}
	return base64.URLEncoding.EncodeToString(code)
}