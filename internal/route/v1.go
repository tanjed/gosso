package route

import (
	"github.com/gorilla/mux"
	"github.com/tanjed/go-sso/internal/handler/auth"
)

func loadV1Routes(router *mux.Router) {
	router.HandleFunc("/token", auth.TokenHandler).Methods("POST")
}