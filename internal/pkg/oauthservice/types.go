package oauthservice

import (
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TokenResponse struct {
	Token string
	Claim jwtmanager.CustomClaims
	Err error
}

type RequiredClaims struct {
	UserId *bson.ObjectID
	ClientId bson.ObjectID
	Scope []string
}