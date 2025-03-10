package jwtmanager

import (
	"errors"
	"fmt"
	"log/slog"
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
		TokenId: bson.NewObjectID(),
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


	oauthTokenPayload := model.NewOauthToken(claims.ClientId, claims.UserId, claims.Scopes, claims.TokenType, claims.ExpAt, claims.IssAt)

	if !oauthTokenPayload.Insert() {
		slog.Error("Unable to store JWT", "error", err)
		return "", errors.New("unable to store jwt")
	}

	return token, nil
}

func ParseToken(tokenStr string) *CustomClaims {
	var customClaims CustomClaims
	app := apiservice.GetApp()
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims, func(token *jwt.Token) (interface{}, error) {
        return []byte(app.Config.JWT_SECRET), nil
    })
	
	if err != nil{
		return nil
	}
	
	if !token.Valid {
		return nil
	}
	return &customClaims
}

func VerifyJwtToken(token *CustomClaims, tokenType string) (model.UserableInterface, error) {	
	
	if !helpers.ContainsInSlice(tokenType, []string{ 
		model.TOKEN_TYPE_USER_ACCESS_TOKEN,  
		model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN})  {
		return nil, &InvalidTokenException{
			Message: "invalid token provided",
			Code: 403,
		}
	}
	user, err := validateCustomClaims(token)

	if err != nil {
		return nil, err
	}
	
	return *user, nil
}

func validateCustomClaims(claims *CustomClaims) (*model.UserableInterface, error) {
    
    oAuthToken := model.GetOAuthTokenById(claims.TokenId)
    if oAuthToken == nil {
        return nil, errors.New("token not found")
    }

	location, _ := time.LoadLocation("Asia/Dhaka")
	if oAuthToken.ExpiredAt.In(location).Before(time.Now().In(location)) {
		return nil, &TokenExpiredException{
			Message: "Token Expired",
			Code: 403,
		}
	}	

	if oAuthToken.Revoked == 1 {
		return nil, &TokenExpiredException{
			Message: "Token Expired",
			Code: 403,
		}
	}
	
	var user model.UserableInterface
	if oAuthToken.TokenType == model.TOKEN_TYPE_USER_ACCESS_TOKEN {
		user = model.GetUserById(oAuthToken.UserId)		
	} else {
		user, _ = model.GetClientById(oAuthToken.ClientId)		
	}

    return &user, nil
}
