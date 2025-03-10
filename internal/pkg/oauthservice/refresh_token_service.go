package oauthservice

import (
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
)


func GetNewRefreshToken(c chan TokenResponse, requiredCalims *RequiredClaims) {
	refreshTokenClaims := jwtmanager.NewJwtClaims(requiredCalims.ClientId,
			&requiredCalims.UserId, requiredCalims.Scope, model.TOKEN_TYPE_USER_REFRESH_TOKEN)

		refreshToken, err := jwtmanager.NewJwtToken(refreshTokenClaims)		

		c <- TokenResponse{
			Token:  refreshToken,
			Claim: *refreshTokenClaims,
			Err: err,
		}
}