package auth

import (
	"net/http"

	"github.com/tanjed/go-sso/internal/customerror"
	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/internal/pkg/oauthservice"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

func passwordGrantHandler(w http.ResponseWriter, r *http.Request) {
	var passwordGrantRequest customtype.PasswordGrantRequest

	if err := passwordGrantRequest.Validated(r.Body); err != nil {
		responsemanager.ResponseUnprocessableEntity(&w, customtype.I{
			"message" : err.Error(),
			"bag" : err.(*customerror.ValidationError).ErrBag,
		})
        return
	}
	
	if _, err := model.GetClientById(passwordGrantRequest.ClientId); err != nil {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid client"})
		return
	}

	user, err := model.AutheticateUser(passwordGrantRequest.MobileNumber, passwordGrantRequest.Password)
	if err != nil  {
		responsemanager.ResponseUnAuthorized(&w, customtype.I{"message" : "invalid user credentials"})
		return
	}


	accessTokenChan := make(chan oauthservice.TokenResponse)
	refreshTokenChan := make(chan oauthservice.TokenResponse)
	requriedClaims := oauthservice.RequiredClaims{
		UserId: user.UserId,
		ClientId: passwordGrantRequest.ClientId,
		Scope: passwordGrantRequest.Scope, 
	}

	go oauthservice.GetNewAccessToken(accessTokenChan, &requriedClaims)
	go oauthservice.GetNewRefreshToken(refreshTokenChan, &requriedClaims)
	
	accessTokenResp := <-accessTokenChan
	refreshTokenResp := <-refreshTokenChan

	if accessTokenResp.Err != nil || refreshTokenResp.Err != nil  {
		responsemanager.ResponseServerError(&w, customtype.I{"message" : "Unable to generate access token"})
		return
	}


	responsemanager.ResponseOK(&w, customtype.I{
		"access_token" : accessTokenResp.Token,
		"refresh_token" : refreshTokenResp.Token,
		"access_token_expired_at" : accessTokenResp.Claim.ExpiresAt,
		"refresh_token_expired_at" : refreshTokenResp.Claim.ExpiresAt,
	})
	
}