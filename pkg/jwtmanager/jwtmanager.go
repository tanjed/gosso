package jwtmanager

import (
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/model"
)

const TOKEN_TYPE_ACCESS_TOKEN = "ACCESS_TOKEN"
const TOKEN_TYPE_REFRESH_TOKEN = "REFRESH_TOKEN"

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

func NewJwtClaims(token_id, clientId, userId string, scopes []string, tokenType string) *CustomClaims{
	expiredAt := time.Now().Add(2 * time.Hour)

	if tokenType == TOKEN_TYPE_REFRESH_TOKEN {
		expiredAt = time.Now().Add((24 * 60) * time.Hour)
	}

	return &CustomClaims{
		TokenId: token_id,
		ClientId: clientId,
		UserId: userId,
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

	oauthTokenPayload := model.NewOauthToken(claims.TokenId, claims.ClientId, claims.UserId, claims.Scopes, 0, claims.TokenType, claims.ExpAt, claims.IssAt, claims.IssAt)

	if !oauthTokenPayload.Insert() {
		slog.Error("Unable to store JWT", "error", err)
		return "", errors.New("unable to store jwt")
	}

	return token, nil
}

func VerifyJwtToken(tokenStr string) bool {
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

	return !validateCustomClaims(customClaims)
}

func validateCustomClaims(claims jwt.Claims) bool{
	return true
}