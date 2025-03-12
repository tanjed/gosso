package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/tanjed/go-sso/internal/handler/customtype"
	"github.com/tanjed/go-sso/internal/model"
	"github.com/tanjed/go-sso/pkg/jwtmanager"
	"github.com/tanjed/go-sso/pkg/responsemanager"
)


func ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {

		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)
		if len(token) <= 0 {
			responsemanager.ResponseUnAuthorized(&w, customtype.M{"message" : "token not found"})
			return
		}

		claims, parsedToken, err := jwtmanager.ParseToken(token)

		if err != nil {
			slog.Error("uable to parse access token", "error", err )
			responsemanager.ResponseUnAuthorized(&w,  customtype.M{"message" : "uable to parse access token"})
			return
		}

		user, err := jwtmanager.VerifyJwtToken(claims, parsedToken)
		
		if err != nil {
			slog.Error("Token Err:", "error", err.Error())
			responsemanager.ResponseUnAuthorized(&w,  customtype.M{"message" : err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), model.AUTH_USER_CONTEXT_KEY, user)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}