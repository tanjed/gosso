package main

import (
	"log"
	"net/http"

	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/internal/route"
)

func main() {
	config.Load()
	db.Init()
	route.Load()

	if err:= http.ListenAndServe(":" + config.AppConfig.Port, nil); err != nil {
		log.Fatalf("Server failed to start %s", err)
	}
}