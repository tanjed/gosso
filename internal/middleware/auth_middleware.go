package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/tanjed/go-sso/internal/model"
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
		user, err := jwtmanager.VerifyJwtToken(parsedToken, parsedToken.TokenType)
		
		if err != nil {
			fmt.Println(err)
			responsemanager.ResponseUnAuthorized(&w, "Invalid token provided")
			return
		}

		var ctx context.Context

		if parsedClient, ok := user.(*model.Client); ok{
			ctx = context.WithValue(r.Context(), model.AUTH_USER_CONTEXT_KEY, parsedClient)
		}

		if parsedUser, ok := user.(*model.User); ok{
			ctx = context.WithValue(r.Context(), model.AUTH_USER_CONTEXT_KEY, parsedUser)
		}
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}