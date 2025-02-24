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

type RefreshTokenRequest struct {
	TokenRequest
	AccessToken string `json:"access_token" validate:"required"`
}

func refreshTokenGrantHandler(w http.ResponseWriter, r *http.Request) {
	var refreshTokenRequest RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
		slog.Error("Unable to process the request", "error", err)
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
        return
	}

	validate := validator.New()

	if err := validate.Struct(refreshTokenRequest); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			responsemanager.ResponseValidationError(&w, errors)
		} else {
			slog.Error("Unable to process the request", "error", err)
			responsemanager.ResponseUnprocessableEntity(&w, "Unable to process the request")
		}
		return
    }
	
	parsedToken := jwtmanager.ParseToken(refreshTokenRequest.AccessToken)

	if parsedToken == nil {
		slog.Error("Uable to parse access token")
		responsemanager.ResponseUnprocessableEntity(&w, "Uable to parse access token")
		return
	}

	var tokenStructType model.TokenableInterface


    switch  parsedToken.TokenType {
    case model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN:
        tokenStructType = &model.ClientAccessToken{}
    case model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN:
        tokenStructType = &model.ClientRefreshToken{}
    case model.TOKEN_TYPE_USER_ACCESS_TOKEN:
        tokenStructType = &model.UserAccessToken{}
    case model.TOKEN_TYPE_USER_REFRESH_TOKEN:
        tokenStructType = &model.UserRefreshToken{}
    }


	oAuthToken := model.GetOAuthTokenById(parsedToken.TokenId, tokenStructType)

	if oAuthToken == nil {
		slog.Error("Invalid token provided")
		responsemanager.ResponseUnprocessableEntity(&w, "Invalid token provided")
		return
	}

	accessTokenClaims := jwtmanager.NewJwtClaims(uuid.New().String(), oAuthToken.GetClientId(),
			nil, oAuthToken.GetScopes(), parsedToken.TokenType)

	accessToken, err := jwtmanager.NewJwtToken(accessTokenClaims)

	if err != nil  {
		responsemanager.ResponseServerError(&w, "Unable to generate access token")
		return
	}

	responsemanager.ResponseOK(&w, map[string]interface{}{
		"access_token" : accessToken,
		"access_token_expired_at" : accessTokenClaims.ExpiresAt,
	})
}