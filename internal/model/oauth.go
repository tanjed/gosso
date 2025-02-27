package model

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
	"go.mongodb.org/mongo-driver/v2/bson"
)
const TOKEN_COLLECTION_NAME = "oauth_tokens"
const TOKEN_TYPE_USER_ACCESS_TOKEN = "USER_ACCESS_TOKEN"
const TOKEN_TYPE_USER_REFRESH_TOKEN = "USER_REFRESH_TOKEN"
const TOKEN_TYPE_CLIENT_ACCESS_TOKEN = "CLIENT_ACCESS_TOKEN"
const TOKEN_TYPE_CLIENT_REFRESH_TOKEN = "CLIENT_REFRESH_TOKEN"

type OauthToken struct {
	TokenId string `bson:"token_id"`
	ClientId string `bson:"client_id"`
	UserId string `bson:"user_id"`
	Scopes []string `bson:"scopes"`
	Revoked int `bson:"revoked"`
	TokenType string `bson:"string"`
	ExpiredAt time.Time `bson:"expired_at"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (t *OauthToken) Insert() bool {
	db := db.InitDB()
	collection := db.Conn.Database(config.AppConfig.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	
	res, err := collection.InsertOne(ctx, t)

	if err != nil {
		slog.Error("Unable to store JWT", "error", err)
		return false
	}

	return res.Acknowledged
}

func (t *OauthToken) InvokeToken() bool{
	db := db.InitDB()
	collection := db.Conn.Database(config.AppConfig.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, bson.D{{"token_id", t.TokenId}}, bson.D{
			{"$set", bson.D{{"revoked", 1}, {"updated_at", time.Now()}},
		},
	})

	if err != nil {
		slog.Error("Unable to invoke token", "error", err)
		return false
	}
	
	return res.Acknowledged
}

func GetOAuthTokenById(tokenId string) *OauthToken{
	db := db.InitDB()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var token OauthToken
	collection := db.Conn.Database(config.AppConfig.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	err := collection.FindOne(ctx, bson.D{{"token_id", tokenId}}).Decode(&token)

	if err != nil {
		slog.Error("Unable to get token", "error", err)
		return nil
	}
	return &token
}

func NewOauthToken(tokenId *string, clientId, userId string, scopes []string, tokenType string, expiredAt time.Time, issusedAt time.Time) *OauthToken {
	if tokenId == nil {
		id := uuid.New().String();
		tokenId = &id
	}

	return &OauthToken{
		TokenId: *tokenId,
		ClientId: clientId,
		UserId: userId,
		Scopes: scopes,
		Revoked: 0,
		TokenType: tokenType,
		ExpiredAt: expiredAt,
		CreatedAt: issusedAt,
		UpdatedAt: issusedAt,
	}
}




