package model

import (
	"context"
	"log/slog"
	"time"

	"github.com/tanjed/go-sso/apiservice"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)


const RESET_PASSWORD_COLLECTION_NAME = "password_resets"

type PasswordReset struct {
	ResetId bson.ObjectID `json:"reset_id" bson:"_id"`
	UserId bson.ObjectID `json:"user_id" bson:"user_id"`
	Token string `json:"token" bson:"token"`
	ExpiredAt time.Time `json:"expired_at" bson:"expired_at"`
	IsValidated int `json:"is_validated" bson:"is_validated"`
	ValidatedAt *time.Time `json:"validated_at" bson:"validated_at"`
	Created_at time.Time `json:"created_at" bson:"created_at"`
}

func NewPasswordReset(userId bson.ObjectID, token string) *PasswordReset {
	return &PasswordReset{
		ResetId: bson.NewObjectID(),
		UserId: userId,
		Token: token,
		ExpiredAt: time.Now().Add(10 * time.Minute),
		ValidatedAt: nil,
		Created_at: time.Now(),
	}
}

func (p *PasswordReset) Insert() bool {
	app := apiservice.GetApp()

	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(RESET_PASSWORD_COLLECTION_NAME)
	res, err := collection.InsertOne(ctx, p)
	if err != nil {
		slog.Error("unable to store otp", "error", err)
		return false
	}

	return res.Acknowledged
}


func GetUserValidResetToken(userId bson.ObjectID) (*PasswordReset, error) {
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(RESET_PASSWORD_COLLECTION_NAME)
	var reset PasswordReset
	err := collection.FindOne(ctx, bson.D{
		{Key: "user_id",  Value: userId},
		{Key: "is_validated", Value: 0},
		{Key: "expired_at", Value: bson.D{{Key: "$gte", Value: time.Now()}}},
	}, options.FindOne().SetSort(bson.D{
		{Key: "created_at", Value: -1},
	})).Decode(&reset)


	if err != nil {
		slog.Error("invalid reset token", "error", err)
		return nil, err
	}

	return &reset, nil
}


func GetResetTokenById(tokenId string) (*PasswordReset, error) {
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(RESET_PASSWORD_COLLECTION_NAME)
	var reset PasswordReset
	err := collection.FindOne(ctx, bson.D{
		{Key: "token_id",  Value: tokenId},
		{Key: "is_validated", Value: 0},
		{Key: "expired_at", Value: bson.D{{Key: "$gte", Value: time.Now()}}},
	}, options.FindOne().SetSort(bson.D{
		{Key: "created_at", Value: -1},
	})).Decode(&reset)


	if err != nil {
		slog.Error("invalid reset token", "error", err)
		return nil, err
	}

	return &reset, nil
}

func (p *PasswordReset) MarkAsValidated() bool {
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(RESET_PASSWORD_COLLECTION_NAME)
	
	res, err := collection.UpdateOne(ctx, bson.D{
		{Key:"reset_id",  Value:p.ResetId},
	}, bson.D{
		{Key:"$set", Value:bson.D{
			{Key:"is_validated",  Value: 1},
			{Key:"validated_at",  Value: time.Now()},
		}},
	})
	if err != nil {
		slog.Error("unable to store otp", "error", err)
		return false
	}
	
	if res.ModifiedCount <= 0 {
		return false
	}
	return true
}