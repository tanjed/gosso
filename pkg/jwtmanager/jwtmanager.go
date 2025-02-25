package jwtmanager

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/helpers"
)

var JWT_SECRET = config.AppConfig.JWT_SECRET


type CustomClaims struct {
	TokenId string
	ClientId string
	UserId string
	Scopes []string
	TokenType string
	ExpAt time.Time
	IssAt time.Time
	jwt.RegisteredClaims
}

func NewJwtClaims(token_id, clientId string, userId *string, scopes []string, tokenType string) *CustomClaims{
	expiredAt := time.Now().Add(2 * time.Hour)

	if tokenType == model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN || tokenType ==  model.TOKEN_TYPE_USER_REFRESH_TOKEN {
		expiredAt = time.Now().Add((24 * 60) * time.Hour)
	}

	return &CustomClaims{
		TokenId: token_id,
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
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(JWT_SECRET))

	if err != nil {
		slog.Error("Unable to generate JWT", "error", err)
		return "", err
	}


	oauthTokenPayload := model.NewOauthToken(&claims.TokenId, claims.ClientId, claims.UserId, claims.Scopes, claims.TokenType, claims.ExpAt, claims.IssAt)

	if !oauthTokenPayload.Insert() {
		slog.Error("Unable to store JWT", "error", err)
		return "", errors.New("unable to store jwt")
	}

	return token, nil
}

func ParseToken(tokenStr string) *CustomClaims {
	var customClaims CustomClaims
	token, err := jwt.ParseWithClaims(tokenStr, &customClaims, func(token *jwt.Token) (interface{}, error) {
        return []byte(JWT_SECRET), nil
    })
	
	if err != nil{
		return nil
	}
	
	if !token.Valid {
		return nil
	}
	fmt.Println(customClaims)
	return &customClaims
}

func VerifyJwtToken(token *CustomClaims, tokenType string) bool {	
	
	if !helpers.ContainsInSlice(tokenType, []string{ 
		model.TOKEN_TYPE_USER_ACCESS_TOKEN,  
		model.TOKEN_TYPE_USER_REFRESH_TOKEN, 
		model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN,  
		model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN})  {
		return false
	}

	return !validateCustomClaims(token)
}

func validateCustomClaims(claims *CustomClaims) bool {
    
    oAuthToken := model.GetOAuthTokenById(claims.TokenId)
    if oAuthToken == nil {
        return false
    }

	if oAuthToken.ExpiredAt.After(time.Now()) {
		return false
	}	

    return false
}
