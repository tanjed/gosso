package jwtmanager

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanjed/go-sso/apiservice"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/helpers"
	"go.mongodb.org/mongo-driver/v2/bson"
)



type CustomClaims struct {
	TokenId bson.ObjectID
	ClientId bson.ObjectID
	UserId bson.ObjectID
	Scopes []string
	TokenType string
	ExpAt time.Time
	IssAt time.Time
	jwt.RegisteredClaims
}

type TokenExpiredException struct {
	Message string
	Code int
}

type InvalidTokenException struct {
	Message string
	Code int
}

func (e TokenExpiredException) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}

func (e InvalidTokenException) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.Message, e.Code)
}


func NewJwtClaims(clientId bson.ObjectID, userId *bson.ObjectID, scopes []string, tokenType string) *CustomClaims{
	expiredAt := time.Now().Add(2 * time.Hour)


	if tokenType == model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN || tokenType ==  model.TOKEN_TYPE_USER_REFRESH_TOKEN {
		expiredAt = time.Now().Add((24 * 60) * time.Hour)
	}

	return &CustomClaims{
		ClientId: clientId,
		UserId: *userId,
		Scopes: scopes,
		TokenType: tokenType,
		ExpAt : expiredAt,
		IssAt : time.Now(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "http://127.0.0.1",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiredAt),
		},
	}
}


func NewJwtToken(claims *CustomClaims) (string, error){
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(apiservice.GetApp().Config.JWT_SECRET))

	if err != nil {
		slog.Error("Unable to generate JWT", "error", err)
		return "", err
	}

	return token, nil
}

func ParseToken(tokenStr string)(*CustomClaims, *jwt.Token, error) {
	var customClaims CustomClaims
	app := apiservice.GetApp()
	t, err := jwt.ParseWithClaims(tokenStr, &customClaims, func(token *jwt.Token) (interface{}, error) {
        return []byte(app.Config.JWT_SECRET), nil
    })
	
	if err != nil{
		return nil, nil, &InvalidTokenException{
			Message: "invalid token",
			Code: http.StatusUnprocessableEntity,
		}
	}
	
	return &customClaims, t, nil
}

func VerifyJwtToken(claims *CustomClaims, token *jwt.Token) (model.UserableInterface, error) {	
	
	if !helpers.ContainsInSlice(claims.TokenType, []string{ 
		model.TOKEN_TYPE_USER_ACCESS_TOKEN,  
		model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN}) || !token.Valid  {
		return nil, &InvalidTokenException{
			Message: "invalid token provided",
			Code: http.StatusUnauthorized,
		}
	}

	user, err := validateCustomClaims(claims)

	if err != nil {
		return nil, err
	}
	
	return *user, nil
}

func validateCustomClaims(claims *CustomClaims) (*model.UserableInterface, error) {
    
    token, err := model.GetOAuthTokenById(claims.TokenId)
    if err != nil {
        return nil, &InvalidTokenException{
			Message: "invalid token provided",
			Code: http.StatusUnauthorized,
		}
    }

	location, _ := time.LoadLocation("Asia/Dhaka")
	if token.ExpiredAt.In(location).Before(time.Now().In(location)) {
		return nil, &TokenExpiredException{
			Message: "token expired",
			Code: http.StatusUnauthorized,
		}
	}	

	if token.Revoked == 1 {
		return nil, &TokenExpiredException{
			Message: "token expired",
			Code: http.StatusUnauthorized,
		}
	}
	
	var user model.UserableInterface
	if token.TokenType == model.TOKEN_TYPE_USER_ACCESS_TOKEN {
		user, err = model.GetUserById(token.UserId)		
	} else {
		user, err = model.GetClientById(token.ClientId)		
	}

	if err != nil {
		return nil, &InvalidTokenException{
			Message: "invalid token",
			Code: http.StatusUnauthorized,
		}
	}

    return &user, nil
}
