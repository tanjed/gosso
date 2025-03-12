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
	Scope []string `bson:"scope"`
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
		Scope: scopes,
		Revoked: 0,
		TokenType: tokenType,
		ExpiredAt: expiredAt,
		CreatedAt: issusedAt,
		UpdatedAt: issusedAt,
	}
}


func (t *OauthToken) Insert() (bson.ObjectID, error) {
	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, t)

	if err != nil {
		slog.Error("Unable to store JWT", "error", err)
		return bson.NilObjectID, err
	}

	return res.InsertedID.(bson.ObjectID), nil
}

func (t *OauthToken) InvokeToken() bool{
	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := collection.UpdateOne(ctx, bson.M{"_id": t.TokenId}, bson.M{
			"$set" : bson.M{"revoked": 1, "updated_at": time.Now()},
		})

	if err != nil {
		slog.Error("unable to invoke token", "error", err)
		return false
	}

	return res.ModifiedCount > 0
}

func GetOAuthTokenById(tokenId bson.ObjectID) (*OauthToken, error){
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	var token OauthToken
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(TOKEN_COLLECTION_NAME)
	err := collection.FindOne(ctx, bson.M{"_id": tokenId}).Decode(&token)

	if err != nil {
		slog.Error("unable to get token", "error", err)
		return nil, err
	}
	return &token, nil
}



