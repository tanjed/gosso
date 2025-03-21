package model

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/apiservice"
	"github.com/tanjed/go-sso/internal/db/redisdb"
	"github.com/tanjed/go-sso/pkg/hashutilities"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const TIME_FORMAT = "2006-01-02 15:04:05.0000"
const CLIENT_COLLECTION_NAME = "clients"

type UserableInterface interface {
	Insert() bool
}
type Client struct {
	ClientId bson.ObjectID 	`json:"client_id" bson:"_id"`
	ClientName string 	`json:"client_name" bson:"client_name"`
	ClientSecret string `json:"client_secret" bson:"client_secret" validate:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type ClientUnAuthorizationCode struct {
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

func NewClient(clientName, clientSecret string) *Client {
	return &Client{
		ClientId: bson.NewObjectID(),
		ClientName: clientName,
		ClientSecret: clientSecret,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (c *Client) Insert() bool {
	app := apiservice.GetApp()
	c.ClientSecret = hashutilities.GenerateHashFromString(c.ClientSecret)
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(CLIENT_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	fmt.Printf("ClientId before insertion: Type=%T, Value=%v\n", c.ClientId, c.ClientId)
	res, err := collection.InsertOne(ctx, c)

	if err != nil {
		slog.Error("Unable to insert client", "error", err)
		return false
	}

	return res.InsertedID != nil
}


func AuthenticateClient(clientId bson.ObjectID, clientSecret string) (*Client, error)  {
	c, err := GetClientById(clientId)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &ClientNotFoundError{
				Message: "client not found",
				Code: http.StatusNotFound,
			}
		}
		return nil, err
	}
	
	
	if !hashutilities.CompareHashWithString(c.ClientSecret, clientSecret) {
		return nil, &ClientUnAuthorizedError{
			Message: "password mismatched",
			Code: http.StatusUnauthorized,
		}
	}

	return c, nil 
}

// func GetClientByClientName(clientName string) *Client {
// 	var client Client
// 	cacheKey := "SSO_CLIENT:" + clientName
// 	if err := redisdb.RedisGetToStruct(cacheKey, &client); err != nil {
// 		if err != redis.Nil {
// 			slog.Error("Unable to get data from redis", "error", err)
// 		}
// 	} else {
// 		return &client
// 	}
// 	app := apiservice.GetApp()

// 	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
// 	defer cancel()
// 	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(CLIENT_COLLECTION_NAME)
// 	err := collection.FindOne(ctx, bson.D{{Key: "client_name", Value: clientName}}).Decode(&client)

// 	if err != nil {
// 		slog.Error("Unable to fetch result", "error", err)
// 		return nil
// 	}
	
// 	if err := redisdb.RedisSetToStruct(cacheKey, &client, (1 * time.Hour)); err != nil {
// 		slog.Error("Unable to set data to redis", "error", err)
// 	}
// 	return &client
// }

func GetClientById(clientId bson.ObjectID) (*Client, error) {
	var client Client
	
	cacheKey := "SSO_CLIENT_BY_CLIENT_ID:" + clientId.Hex()

	if err := redisdb.RedisGetToStruct(cacheKey, &client); err != nil {
		if err != redis.Nil {
			slog.Error("unable to get data from redis", "error", err)
		}
	} else {
		
		return &client, nil
	}

	app := apiservice.GetApp()

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(CLIENT_COLLECTION_NAME)
	err := collection.FindOne(ctx, bson.M{"_id": clientId}).Decode(&client)

	if err != nil {
		slog.Error("unable to fetch result", "error", err)
		return nil, err
	}
	
	if err := redisdb.RedisSetToStruct(cacheKey, &client, (1 * time.Hour)); err != nil {
		slog.Error("unable to set data to redis", "error", err)
	}
	return &client, nil
}

// func ClientHasValidSession(clientId string, redirectUri string) (*clientUnAuthorizationCode, error) {
// 	db := db.InitDB()
// 	var clientAuthorizationCode clientUnAuthorizationCode
// 	err := db.Conn.Query("SELECT client_id, client_code, redirect_uri, generated_at, expired_at FROM client_authorization_codes WHERE client_id = ? AND redirect_uri = ? AND expired_at > ?", clientId, redirectUri, time.Now()).
// 	Scan(&clientAuthorizationCode.ClientId, &clientAuthorizationCode.ClientCode, &clientAuthorizationCode.RedirectUri, &clientAuthorizationCode.GeneratedAt, &clientAuthorizationCode.ExpiredAt)

// 	if err != nil {
// 		log.Println("Unable to fetch result", err)
// 		return nil, err
// 	}

// 	return &clientAuthorizationCode, nil
// }

// func NewClientAuthorizationCode(clientId string, redirectUri string) (*clientUnAuthorizationCode, error){
// 	clientUnAuthorizationCode := clientUnAuthorizationCode {
// 		ClientId: clientId,
// 		ClientCode: generateAuthorizationCode(),
// 		RedirectUri: redirectUri,
// 		GeneratedAt: time.Now(),
// 		ExpiredAt: time.Now().Add(30 * time.Second),
// 	}
// 	err := insertClientAuthorizationCode(clientUnAuthorizationCode)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &clientUnAuthorizationCode, nil
// }

// func insertClientAuthorizationCode(c clientUnAuthorizationCode) error {
// 	db := db.InitDB()

// 	return db.Conn.Query("INSERT INTO client_authorization_codes (client_id, client_code, redirect_uri, generated_at, expired_at) VALUES (?,?,?,?,?)", c.ClientId, c.ClientCode, c.RedirectUri, c.GeneratedAt, c.ExpiredAt).Exec()
// }

// func generateAuthorizationCode() string {
// 	code := make([]byte, 32)
// 	_, err := rand.Read(code)
// 	if err != nil {
// 		log.Fatal("Failed to generate authorization code:", err)
// 	}
// 	return base64.URLEncoding.EncodeToString(code)
// }


