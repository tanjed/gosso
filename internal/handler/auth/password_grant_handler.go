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

type PasswordGrantRequest struct {
	TokenRequest
	ClientName string 	`json:"client_name" validate:"required"`
	ClientSecret string `json:"client_secret" validate:"required"`
	Scope []string 		`json:"scope" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
	Password string 	`json:"password" validate:"required"`
}

func passwordGrantHandler(w http.ResponseWriter, r *http.Request) {
	var passwordGrantRequest PasswordGrantRequest
	err := json.NewDecoder(r.Body).Decode(&passwordGrantRequest); 
	
	if err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()
	err = validate.Struct(passwordGrantRequest)

	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }
	
	client, err := model.AuthenticateClient(passwordGrantRequest.ClientName, passwordGrantRequest.ClientSecret);
	if err != nil {
		responsemanager.ResponseUnAuthorized(&w, "Invalid client credentials")
		return
	}

	user, err := model.AutheticateUser(passwordGrantRequest.MobileNumber, passwordGrantRequest.Password)
	if err != nil  {
		responsemanager.ResponseUnAuthorized(&w, "Invalid user credentials")
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
			&user.UserId, passwordGrantRequest.Scope, model.TOKEN_TYPE_USER_ACCESS_TOKEN)

			accessToken, err := jwtmanager.NewJwtToken(accessTokenClaims)

		accessTokenChan	<- tokenResponse{
			Token: accessToken,
			Claim: *accessTokenClaims,
			Err: err,
		}
	}()


	go func ()  {
		refreshTokenClaims := jwtmanager.NewJwtClaims(uuid.New().String(), client.ClientId,
			&user.UserId, passwordGrantRequest.Scope, model.TOKEN_TYPE_USER_REFRESH_TOKEN)

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