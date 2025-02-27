package model

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/pkg/hashutilities"
	"go.mongodb.org/mongo-driver/v2/bson"
)


const USER_COLLECTION_NAME = "users"
type User struct {
	UserId string `bson:"user_id"`
	FirstName string `bson:"first_name"`
	LastName string `bson:"last_name"`
	MobileNumber string `bson:"mobile_number"`
	Password string `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`

}

type UserNotFound struct {
	Message string
	Code int
}

type UserUnauthorized struct {
	Message string
	Code int
}

func (e UserNotFound) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

func (e UserUnauthorized) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

func (u *User) Insert() bool {
	db := db.InitDB()
	hashedPassword := hashutilities.GenerateHashFromString(u.Password)
	collection := db.Conn.Database(config.AppConfig.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, User{
		UserId: u.UserId,
		FirstName: u.FirstName,
		LastName: u.LastName,
		MobileNumber: u.MobileNumber,
		Password: hashedPassword,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})

	if err != nil {
		slog.Error("Unable to insert user", "error", err)
		return false
	}

	return res.Acknowledged
}

func AutheticateUser(mobileNumber string, password string) (*User, error) {

	u := GetUserByMobileNumber(mobileNumber)
	
	if u == nil {
		return nil, &UserNotFound{
			Message: "User not found",
			Code: http.StatusNotFound,
		}
	}
	user := *u
	if !hashutilities.CompareHashWithString(user.Password, password) {
		return nil, &UserUnauthorized{
			Message: "User unauthorized",
			Code: http.StatusUnauthorized,
		}
	}
	return u, nil
}

func GetUserByMobileNumber(mobileNumber string) *User {
	var user User
	cacheKey := "SSO_USER:" + mobileNumber

	if err := db.RedisGetToStruct(cacheKey, &user); err != nil {
		if err != redis.Nil {
			slog.Error("Unable to get data from redis", "error", err)
		}
	} else {
		
		return &user
	}
	
	dbConn := db.InitDB()
	collection := dbConn.Conn.Database(config.AppConfig.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.D{{"mobile_number", mobileNumber}}).Decode(&user)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil
	}

	if err := db.RedisSetToStruct(cacheKey, &user, (1 * time.Hour)); err != nil {
		slog.Error("Unable to set data to redis", "error", err)
	}

	return &user
}


func GetUserById(userId string) *User {
	var user User
	cacheKey := "SSO_USER_BY_ID:" + userId

	if err := db.RedisGetToStruct(cacheKey, &user); err != nil {
		if err != redis.Nil {
			slog.Error("Unable to get data from redis", "error", err)
		}
	} else {
		
		return &user
	}
	
	dbConn := db.InitDB()
	collection := dbConn.Conn.Database(config.AppConfig.DB_NAME).Collection(USER_COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.D{{"user_id", userId}}).Decode(&user)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil
	}

	if err := db.RedisSetToStruct(cacheKey, &user, (1 * time.Hour)); err != nil {
		slog.Error("Unable to set data to redis", "error", err)
	}

	return &user
}

func NewUser(firstName, lastName, mobileNumber, password string) *User {
	return &User{
		UserId: uuid.New().String(),
		FirstName : firstName,
		LastName : lastName,
		MobileNumber : mobileNumber,
		Password : password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}