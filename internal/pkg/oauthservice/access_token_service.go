package oauthservice

import (
	"log/slog"

	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
)


func GetOAuthToken(c chan TokenResponse, tokenType string, requiredCalims *RequiredClaims) {
	claims := jwtmanager.NewJwtClaims(requiredCalims.ClientId, requiredCalims.UserId, requiredCalims.Scope, tokenType)
	oauthTokenPayload := model.NewOauthToken(claims.ClientId, claims.UserId, claims.Scopes, claims.TokenType, claims.ExpAt, claims.IssAt)

	tokenId, err := oauthTokenPayload.Insert();
	claims.TokenId = tokenId

	if err != nil {
		slog.Error("Unable to store JWT", "error", err)
		c <- TokenResponse{
			Err: err,
		}
		return 
	}

	accessToken, err := jwtmanager.NewJwtToken(claims)
		
	c	<- TokenResponse{
		Token: accessToken,
		Claim: *claims,
		Err: err,
	}
}