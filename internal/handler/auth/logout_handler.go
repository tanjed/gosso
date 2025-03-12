package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

type LogoutRequest struct {
	customtype.TokenRequest
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var logoutRequest LogoutRequest

	json.NewDecoder(r.Body).Decode(&logoutRequest)

	claims, parsedToken, err := jwtmanager.ParseToken(strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1))

	if err != nil {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : err.Error()})
		return
	}

	if !parsedToken.Valid {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid token"})
		return
	}
	
	token, err := model.GetOAuthTokenById(claims.TokenId)
	
	if err != nil || token.Revoked == 1 {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "invalid token provided"})
		return
	}
	
	if isInvoked := token.InvokeToken(); !isInvoked {
		responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "unable to invoke token"})
		return
	}

	responsemanager.ResponseOK(&w, customtype.M{"message" : "token invoked"})

}