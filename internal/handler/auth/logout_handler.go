package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type LogoutRequest struct {
	TokenRequest
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var logoutRequest LogoutRequest

	json.NewDecoder(r.Body).Decode(&logoutRequest)
	token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
	claims := jwtmanager.ParseToken(token)
	if claims == nil {
		responsemanager.ResponseUnprocessableEntity(&w, "Invalid token provided")
		return
	}
	var tokenStructType model.TokenableInterface
	switch claims.TokenType {
    case model.TOKEN_TYPE_CLIENT_ACCESS_TOKEN:
        tokenStructType = &model.ClientAccessToken{}
    case model.TOKEN_TYPE_CLIENT_REFRESH_TOKEN:
        tokenStructType = &model.ClientRefreshToken{}
    case model.TOKEN_TYPE_USER_ACCESS_TOKEN:
        tokenStructType = &model.UserAccessToken{}
    case model.TOKEN_TYPE_USER_REFRESH_TOKEN:
        tokenStructType = &model.UserRefreshToken{}
    }


	oAuthToken := model.GetOAuthTokenById(claims.TokenId, tokenStructType)

	if oAuthToken == nil {
		responsemanager.ResponseUnprocessableEntity(&w, "Invalid token provided")
		return
	}
	
	if isInvoked := oAuthToken.InvokeToken(); !isInvoked {
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to invoke token")
		return
	}

	responsemanager.ResponseOK(&w, "Token invoked")

}