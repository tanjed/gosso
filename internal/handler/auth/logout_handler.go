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

	oAuthToken := model.GetOAuthTokenById(claims.ID)
	isInvoked := false

	if logoutRequest.GrantType == GRANT_TYPE_CLIENT_CREDENTIALS {
		isInvoked = oAuthToken.InvokeClient()
	} else if logoutRequest.GrantType == GRANT_TYPE_PASSWORD {
		isInvoked = oAuthToken.InvokeUser()
	} else {
		isInvoked = oAuthToken.InvokeToken()
	}

	if !isInvoked {
		responsemanager.ResponseUnprocessableEntity(&w, "Unable to invoke token")
		return
	}

	responsemanager.ResponseOK(&w, "Token invoked")

}