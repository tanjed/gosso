package model

import (
	"context"
	"log/slog"
	"time"

	"github.com/tanjed/go-sso/apiservice"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)
const OTP_COLLECTION_NAME = "otps"

const OTP_TYPE_PASSWORD_RESET = "password_reset"

const MAX_RETRY_COUNT = 3

type Otp struct {
	OtpId bson.ObjectID `bson:"_id"`
	UserId bson.ObjectID `bson:"user_id"`
	Otp string `bson:"otp"`
	OtpType string `bson:"otp_type"`
	IsValidated int `bson:"is_validated"`
	ExpiredAt time.Time `bson:"expired_at"`
	ValidatedAt *time.Time `bson:"validated_at"`
	CreatedAt time.Time `bson:"created_at"`
}

func NewOtp(userId bson.ObjectID, otp, otpType string) *Otp{	
	return &Otp{
		OtpId: bson.NewObjectID(),
		UserId: userId,
		Otp: otp,
		OtpType: otpType,
		IsValidated: 0,
		ExpiredAt:  time.Now().Add(5 * time.Minute),
		ValidatedAt: nil,
		CreatedAt: time.Now(),
	}
}

func (o *Otp) Insert() bool {
	app := apiservice.GetApp()

	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(OTP_COLLECTION_NAME)
	res, err := collection.InsertOne(ctx, o)
	if err != nil {
		slog.Error("unable to store otp", "error", err)
		return false
	}

	return res.Acknowledged
}

func (o *Otp) MarkAsValidated() bool {
	app := apiservice.GetApp()

	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(OTP_COLLECTION_NAME)
	res, err := collection.UpdateOne(ctx, bson.D{
		{ "otp_id",  o.OtpId},
	}, bson.D{
		{"$set", bson.D{
			{ "is_validated",  1},
			{ "validated_at",  time.Now()},
		}},
	})
	if err != nil {
		slog.Error("unable to mark otp as validated", "error", err)
		return false
	}

	return res.Acknowledged
}

func GetUserSentOtpCount(userId bson.ObjectID, otpType string) (int64, error) {
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(OTP_COLLECTION_NAME)
	now := time.Now()
	todayStartTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEndTime := todayStartTime.Add(24 * time.Hour)

	count, err := collection.CountDocuments(ctx, bson.D{
		{Key: "expired_at",Value: bson.D{{Key: "$gte",Value:  time.Now()}}},
		{Key: "created_at",Value: bson.D{{Key: "$gte",Value:  todayStartTime}, {Key: "$lte",Value: todayEndTime}}},
		{Key: "user_id",Value: userId},	
		{Key: "is_validated",Value: 0},
		{Key: "otp_type",Value: otpType},
	})

	if err != nil {
		slog.Error("unable to store otp", "error", err)
		return 0, err
	}

	return count, nil
}


func GetUserValidOtp(userId bson.ObjectID, otpType string) (*Otp, error) {
	app := apiservice.GetApp()
	ctx, cancel := context.WithTimeout(context.Background(), (5 * time.Second))
	defer cancel()

	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(OTP_COLLECTION_NAME)
	var otp Otp
	err := collection.FindOne(ctx, bson.D{
		{Key: "user_id",  Value: userId},
		{Key:"otp_type",  Value:otpType},
		{Key: "is_validated", Value:0},
		{Key: "expired_at", Value:bson.D{{Key: "$gte", Value: time.Now()}}},
	}, options.FindOne().SetSort(bson.D{
		{Key: "created_at", Value: -1},
	})).Decode(&otp)


	if err != nil {
		slog.Error("valid otp not found", "error", err)
		return nil, err
	}

	return &otp, nil
}