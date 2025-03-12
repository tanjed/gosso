package auth

import (
	"log/slog"
	"net/http"

	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/pkg/oauthservice"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)


func refreshTokenGrantHandler(w http.ResponseWriter, r *http.Request) {
	var refreshTokenRequest customtype.RefreshTokenRequest

	if err := refreshTokenRequest.Validated(r.Body); err != nil {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.I{
			"message" : err.Error(),
			"bag" : err.(*customerror.ValidationError).ErrBag,
		})
        return
	}
	
	claims, parsedToken, err := jwtmanager.ParseToken(refreshTokenRequest.AccessToken)

	if err != nil {
		slog.Error("uable to parse access token", "error", err )
		responsemanager.ResponseUnAuthorized(&w, "uable to parse access token")
		return
	}

	if !parsedToken.Valid {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid refresh token"})
		return
	}

	token, err := model.GetOAuthTokenById(claims.TokenId)

	if err != nil || token.Revoked == 1 {
		slog.Error("invalid token provided", "error", err)
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid refresh token"})
		return
	}

	if isInvoked := token.InvokeToken(); !isInvoked {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "unable to process refresh token"})
		return
	}

	accessTokenChan := make(chan oauthservice.TokenResponse)
	refreshTokenChan := make(chan oauthservice.TokenResponse)
	requiredClaims := oauthservice.RequiredClaims{
		UserId: &token.UserId,
		ClientId: token.ClientId,
		Scope: token.Scope, 
	}

	go oauthservice.GetOAuthToken(accessTokenChan, model.TOKEN_TYPE_USER_ACCESS_TOKEN, &requiredClaims)
	go oauthservice.GetOAuthToken(refreshTokenChan, model.TOKEN_TYPE_USER_REFRESH_TOKEN, &requiredClaims)

	accessTokenResp := <-accessTokenChan
	refreshTokenResp := <-refreshTokenChan

	if accessTokenResp.Err != nil || refreshTokenResp.Err != nil  {
		responsemanager.ResponseServerError(&w, customtype.I{"message" : "unable to generate token"})
		return
	}

	responsemanager.ResponseOK(&w, customtype.I{
		"access_token" : accessTokenResp.Token,
		"refresh_token" : refreshTokenResp.Token,
		"access_token_expired_at" : accessTokenResp.Claim.ExpiresAt,
		"refresh_token_expired_at" : refreshTokenResp.Claim.ExpiresAt,
	})
}