package route

import (
	"github.com/gorilla/mux"
	"github.com/tanjed/go-sso/internal/handler/auth"
	"github.com/tanjed/go-sso/internal/handler/profile"
	"github.com/tanjed/go-sso/internal/middleware"
)

func loadV1Routes(router *mux.Router) {
	

	authRouter := router.NewRoute().Subrouter()
	authRouter.Use(middleware.ValidateToken)

	loadPublicRoutes(router)
	loadPrivateRoutes(authRouter)
}

func loadPublicRoutes(router *mux.Router) {
	router.HandleFunc("/token", auth.TokenHandler).Methods("POST")
	router.HandleFunc("/register", auth.UserRegisterHandler).Methods("POST")
}

func loadPrivateRoutes(router *mux.Router) {
	router.HandleFunc("/invoke", auth.LogoutHandler).Methods("POST")
	router.HandleFunc("/profile", profile.UserProfileHandler).Methods("GET")
}