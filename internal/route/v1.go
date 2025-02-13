package route

import (
	"github.com/gorilla/mux"
	"github.com/tanjed/go-sso/internal/handler/auth"
)

func loadV1Routes(router *mux.Router) {
	router.HandleFunc("/authorize", auth.AuthorizeHandler).Methods("POST")
	router.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
	router.HandleFunc("/token", auth.TokenHandler).Methods("POST")
	router.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
}