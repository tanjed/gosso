package model

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/apiservice"
	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/db/redisdb"
	"github.com/tanjed/go-sso/internal/handler/customtype"

	"github.com/tanjed/go-sso/pkg/hashutilities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const USER_COLLECTION_NAME = "users"
type User struct {
	UserId bson.ObjectID `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName string `bson:"last_name"`
	MobileNumber string `bson:"mobile_number"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	Address *string `bson:"address"`
	Nid *string `bson:"nid"`
	Passport *string `bson:"passport"`
	ClientId bson.ObjectID `bson:"client_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

}

func NewUser(d customtype.UserRegisterRequest) *User {
	return &User{
		UserId: bson.NewObjectID(),
		FirstName : d.FirstName,
		LastName : d.LastName,
		MobileNumber : d.MobileNumber,
		Email: d.Email,
		Password : d.Password,
		Address : &d.Address,
		Nid : &d.Nid,
		Passport : &d.Passport,
		ClientId : d.ClientId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) Insert() bool {
	app := apiservice.GetApp()
	u.Password = hashutilities.GenerateHashFromString(u.Password)
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, u)

	if err != nil {
		slog.Error("Unable to insert user", "error", err)
		return false
	}

	return res.Acknowledged
}

func AutheticateUser(mobileNumber string, password string) (*User, error) {

	u, err := GetUserByMobileNumber(mobileNumber)

	if err != nil {
		return nil, &customerror.UserNotFoundError{
			ErrMessage: "User not found",
			ErrCode: http.StatusNotFound,
		}
	}
	
	if !hashutilities.CompareHashWithString(u.Password, password) {
		return nil, &customerror.UserUnauthorizedError{
			ErrMessage: "User unauthorized",
			ErrCode: http.StatusUnauthorized,
		}
	}
	return u, nil
}

func GetUserByMobileNumber(mobileNumber string) (*User, error) {
	var user User
	cacheKey := "SSO_USER:" + mobileNumber

	if err := redisdb.RedisGetToStruct(cacheKey, &user); err != nil {
		if err != redis.Nil {
			slog.Error("Unable to get data from redis", "error", err)
		}
	} else {

		return &user, nil
	}

	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"mobile_number": mobileNumber}).Decode(&user)

	if err != nil {
		slog.Error("Unable to fetch result from db", "error", err)
		return nil, err
	}

	if err := redisdb.RedisSetToStruct(cacheKey, &user, (1 * time.Second)); err != nil {
		slog.Error("Unable to set data to redis", "error", err)
	}

	return &user, nil
}

func GetUserById(userId bson.ObjectID) *User {
	var user User
	cacheKey := "SSO_USER_BY_ID:" + userId.Hex()

	if err := redisdb.RedisGetToStruct(cacheKey, &user); err != nil {
		if err != redis.Nil {
			slog.Error("Unable to get data from redis", "error", err)
		}
	} else {

		return &user
	}

	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.D{{Key:"user_id", Value: userId}}).Decode(&user)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil
	}

	if err := redisdb.RedisSetToStruct(cacheKey, &user, (1 * time.Hour)); err != nil {
		slog.Error("Unable to set data to redis", "error", err)
	}

	return &user
}

func (u *User) UpdatePassword(newPassword string) (bool){
	hashedPassword := hashutilities.GenerateHashFromString(newPassword)
	app := apiservice.GetApp()
	collection := app.DB.Conn.Database(app.Config.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, bson.D{{Key: "user_id", Value: u.UserId}}, bson.D{
		{Key: "$set", Value: bson.D{{Key:"password", Value: hashedPassword}}},
	})

	if err != nil {
		log.Println("Unable to update", err)
		return false
	}

	return res.Acknowledged
}