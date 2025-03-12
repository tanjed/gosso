package auth

import (
	"net/http"

	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/pkg/oauthservice"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)


func clientCredentialsGrantHandler(w http.ResponseWriter, r *http.Request) {
	var clientCredentialsGrantRequest customtype.ClientCredentialsGrantRequest

	if err := clientCredentialsGrantRequest.Validated(r.Body); err != nil {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.I{
			"message" : err.Error(),
			"bag" : err.(*customerror.ValidationError).ErrBag,
		})
        return
	}
	
	client, err := model.AuthenticateClient(clientCredentialsGrantRequest.ClientId, clientCredentialsGrantRequest.ClientSecret);
	if err != nil {
		switch err.(type) {
		case model.ClientNotFoundError :
			responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid client credentials"})
		default:
			responsemanager.ResponseServerError(&w, customtype.M{"message" : "something went wrong"})
		}
		return
	}

	
	accessTokenChan := make(chan oauthservice.TokenResponse)
	refreshTokenChan := make(chan oauthservice.TokenResponse)
	requiredClaims := oauthservice.RequiredClaims{
		UserId: nil,
		ClientId: client.ClientId,
		Scope: clientCredentialsGrantRequest.Scope,
	}

	go oauthservice.GetOAuthToken(accessTokenChan, model.TOKEN_TYPE_USER_ACCESS_TOKEN, &requiredClaims)
	go oauthservice.GetOAuthToken(refreshTokenChan, model.TOKEN_TYPE_USER_REFRESH_TOKEN, &requiredClaims)

	
	accessTokenResp := <-accessTokenChan
	refreshTokenResp := <-refreshTokenChan

	if accessTokenResp.Err != nil || refreshTokenResp.Err != nil  {
		responsemanager.ResponseServerError(&w, customtype.M{"message" : "unable to generate access token"})
		return
	}


	responsemanager.ResponseOK(&w, customtype.I{
		"access_token" : accessTokenResp.Token,
		"refresh_token" : refreshTokenResp.Token,
		"access_token_expired_at" : accessTokenResp.Claim.ExpiresAt,
		"refresh_token_expired_at" : refreshTokenResp.Claim.ExpiresAt,
	})
	
}