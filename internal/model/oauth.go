package model

import (
	"context"
	"log/slog"
	"time"

	"github.com/tanjed/go-sso/apiservice"
	"go.mongodb.org/mongo-driver/v2/bson"
)
const TOKEN_COLLECTION_NAME = "oauth_tokens"
const TOKEN_TYPE_USER_ACCESS_TOKEN = "USER_ACCESS_TOKEN"
const TOKEN_TYPE_USER_REFRESH_TOKEN = "USER_REFRESH_TOKEN"
const TOKEN_TYPE_CLIENT_ACCESS_TOKEN = "CLIENT_ACCESS_TOKEN"
const TOKEN_TYPE_CLIENT_REFRESH_TOKEN = "CLIENT_REFRESH_TOKEN"

type OauthToken struct {
	TokenId bson.ObjectID `bson:"_id"`
	ClientId bson.ObjectID `bson:"client_id"`
	UserId bson.ObjectID `bson:"user_id"`
	Scopes []string `bson:"scopes"`
	Revoked int `bson:"revoked"`
	TokenType string `bson:"string"`
	ExpiredAt time.Time `bson:"expired_at"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewOauthToken(clientId bson.ObjectID, userId bson.ObjectID, scopes []string, tokenType string, expiredAt time.Time, issusedAt time.Time) *OauthToken {
	// if tokenId == nil {
	// 	id := uuid.New().String();
	// 	tokenId = &id
	// }

	return &OauthToken{
		TokenId: bson.NewObjectID(),
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


func (t *OauthToken) Insert() bool {
	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
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
	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
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

func GetOAuthTokenById(tokenId bson.ObjectID) *OauthToken{
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var token OauthToken
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	err := collection.FindOne(ctx, bson.M{"_id": tokenId}).Decode(&token)

	if err != nil {
		slog.Error("Unable to get token", "error", err)
		return nil
	}
	return &token
}



