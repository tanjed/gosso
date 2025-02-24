package jwtmanager

import (
	"errors"
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

	var tokenStructType model.TokenableInterface

    switch claims.TokenType {
    case model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN:
        tokenStructType = &model.ClientAccessToken{}
    case model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN:
        tokenStructType = &model.ClientRefreshToken{}
    case model.TOKEN_TYPE_USER_ACCESS_TOKEN:
        tokenStructType = &model.UserAccessToken{}
    case model.TOKEN_TYPE_USER_REFRESH_TOKEN:
        tokenStructType = &model.UserRefreshToken{}
    }


	oauthTokenPayload := NewOauthToken(tokenStructType, claims)

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
	
	return &customClaims
}

func VerifyJwtToken(tokenStr string, tokenType string) bool {	
	
	if !helpers.ContainsInSlice(tokenType, []string{ 
		model.TOKEN_TYPE_USER_ACCESS_TOKEN,  
		model.TOKEN_TYPE_USER_REFRESH_TOKEN, 
		model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN,  
		model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN})  {
		return false
	}

	var customClaims CustomClaims
	token, err := jwt.ParseWithClaims(tokenStr, customClaims, func(token *jwt.Token) (interface{}, error) {
        return JWT_SECRET, nil
    })

	if err != nil{
		return false
	}

	if !token.Valid {
		return false
	}

	return !validateCustomClaims(customClaims, tokenType)
}

func validateCustomClaims(claims CustomClaims, tokenType string) bool {
    var tokenStructType model.TokenableInterface

    switch tokenType {
    case model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN:
        tokenStructType = &model.ClientAccessToken{}
    case model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN:
        tokenStructType = &model.ClientRefreshToken{}
    case model.TOKEN_TYPE_USER_ACCESS_TOKEN:
        tokenStructType = &model.UserAccessToken{}
    case model.TOKEN_TYPE_USER_REFRESH_TOKEN:
        tokenStructType = &model.UserRefreshToken{}
    }


    oAuthToken := model.GetOAuthTokenById(claims.TokenId, tokenStructType)
    if oAuthToken == nil {
        return false
    }

	if oAuthToken.GetExpiry().After(time.Now()) {
		return false
	}	

    return false
}

func NewOauthToken(tokenableStruct model.TokenableInterface, claims *CustomClaims) model.TokenableInterface{
	if token, ok := tokenableStruct.(*model.ClientAccessToken); ok {
		token.TokenId = claims.TokenId
		token.ClientId = claims.ClientId
		token.Scopes = claims.Scopes
		token.Revoked = 0
		token.ExpiredAt = claims.ExpAt
		token.CreatedAt = claims.IssAt
		token.UpdatedAt = claims.IssAt
		return token
	} else if token, ok := tokenableStruct.(*model.ClientRefreshToken); ok {
		token.TokenId = claims.TokenId
		token.ClientId = claims.ClientId
		token.Scopes = claims.Scopes
		token.Revoked = 0
		token.ExpiredAt = claims.ExpAt
		token.CreatedAt = claims.IssAt
		token.UpdatedAt = claims.IssAt
		return token
	} else if token, ok := tokenableStruct.(*model.UserAccessToken); ok {
		token.TokenId = claims.TokenId
		token.ClientId = claims.ClientId
		token.UserId = claims.UserId
		token.Scopes = claims.Scopes
		token.Revoked = 0
		token.ExpiredAt = claims.ExpAt
		token.CreatedAt = claims.IssAt
		token.UpdatedAt = claims.IssAt
		return token
		
	} else if token, ok := tokenableStruct.(*model.UserRefreshToken); ok {
		token.TokenId = claims.TokenId
		token.ClientId = claims.ClientId
		token.UserId = claims.UserId
		token.Scopes = claims.Scopes
		token.Revoked = 0
		token.ExpiredAt = claims.ExpAt
		token.CreatedAt = claims.IssAt
		token.UpdatedAt = claims.IssAt
		return token
	} else {
		return nil
	}
}