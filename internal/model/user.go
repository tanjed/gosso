package model

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/pkg/hashutilities"
)

type User struct {
	UserId string
	FirstName string
	LastName string
	MobileNumber string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time

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


func AutheticateUser(mobileNumber string, password string) (*User, error) {

	u := getUserByMobileNumber(mobileNumber)
	
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

func getUserByMobileNumber(mobileNumber string) *User {
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
	
	err := dbConn.Conn.Query("SELECT user_id, first_name, last_name, mobile_number, password, created_at, updated_at FROM users WHERE mobile_number = ?", mobileNumber).Scan(&user.UserId,&user.FirstName, &user.LastName, &user.MobileNumber, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Println("Unable to fetch result", err)
		return nil
	}

	if err := db.RedisSetToStruct(cacheKey, &user, (1 * time.Hour)); err != nil {
		slog.Error("Unable to set data to redis", "error", err)
	}

	return &user
}