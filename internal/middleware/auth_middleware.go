package middleware

import (
	"net/http"
	"strings"

	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)

func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		if len(token) <= 0 {
			responsemanager.ResponseUnAuthorized(&w, "Token not provided")
			return
		}
		parsedToken := jwtmanager.ParseToken(token)
		if !jwtmanager.VerifyJwtToken(token, parsedToken.TokenType) {
			responsemanager.ResponseUnAuthorized(&w, "Invalid token provided")
			return
		}
		next.ServeHTTP(w, r)
	})
}