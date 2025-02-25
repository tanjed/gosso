package auth

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(claims)
	if claims == nil {
		responsemanager.ResponseUnprocessableEntity(&w, "Invalid token provided")
		return
	}
	
	fmt.Println(claims.TokenId)
	oAuthToken := model.GetOAuthTokenById(claims.TokenId)
	
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