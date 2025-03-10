package oauthservice

import (
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
)


func GetNewAccessToken(c chan TokenResponse, requiredCalims *RequiredClaims) {
	accessTokenClaims := jwtmanager.NewJwtClaims(requiredCalims.ClientId,
			&requiredCalims.UserId, requiredCalims.Scope, model.TOKEN_TYPE_USER_ACCESS_TOKEN)

			accessToken, err := jwtmanager.NewJwtToken(accessTokenClaims)

		c	<- TokenResponse{
			Token: accessToken,
			Claim: *accessTokenClaims,
			Err: err,
		}
}