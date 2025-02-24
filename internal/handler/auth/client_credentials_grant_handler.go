package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type ClientCredentialsGrantRequest struct {
	TokenRequest
	ClientName string 	`json:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Scope []string 		`json:"scope" validate:"required"`
}

func clientCredentialsGrantHandler(w http.ResponseWriter, r *http.Request) {
	var clientCredentialsGrantRequest ClientCredentialsGrantRequest
	err := json.NewDecoder(r.Body).Decode(&clientCredentialsGrantRequest); 
	
	if err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()
	err = validate.Struct(clientCredentialsGrantRequest)

	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }
	
	client, err := model.AuthenticateClient(clientCredentialsGrantRequest.ClientName, clientCredentialsGrantRequest.ClientSecret);
	if err != nil {
		responsemanager.ResponseUnAuthorized(&w, "Invalid client credentials")
		return
	}

	type tokenResponse struct {
		Token string
		Claim jwtmanager.CustomClaims
		Err error
	}

	accessTokenChan := make(chan tokenResponse)
	refreshTokenChan := make(chan tokenResponse)

	go func ()  {
		accessTokenClaims := jwtmanager.NewJwtClaims(uuid.New().String(), client.ClientId,
			nil, clientCredentialsGrantRequest.Scope, model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN)

			accessToken, err := jwtmanager.NewJwtToken(accessTokenClaims)

		accessTokenChan	<- tokenResponse{
			Token: accessToken,
			Claim: *accessTokenClaims,
			Err: err,
		}
	}()


	go func ()  {
		refreshTokenClaims := jwtmanager.NewJwtClaims(uuid.New().String(), client.ClientId,
			nil, clientCredentialsGrantRequest.Scope, model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN)

		refreshToken, err := jwtmanager.NewJwtToken(refreshTokenClaims)		

		refreshTokenChan <- tokenResponse{
			Token:  refreshToken,
			Claim: *refreshTokenClaims,
			Err: err,
		}
	}()
	
	accessTokenResp := <-accessTokenChan
	refreshTokenResp := <-refreshTokenChan

	if accessTokenResp.Err != nil || refreshTokenResp.Err != nil  {
		responsemanager.ResponseServerError(&w, "Unable to generate access token")
		return
	}


	responsemanager.ResponseOK(&w, map[string]interface{}{
		"access_token" : accessTokenResp.Token,
		"refresh_token" : refreshTokenResp.Token,
		"access_token_expired_at" : accessTokenResp.Claim.ExpiresAt,
		"refresh_token_expired_at" : refreshTokenResp.Claim.ExpiresAt,
	})
	
}