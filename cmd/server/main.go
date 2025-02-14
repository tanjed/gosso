package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db"
	"github.com/tanjed/go-sso/internal/route"
)

func main() {
	config.Load()

	db := db.Init()
	db.Close()
	r := route.Load()
	server := http.Server{
		Addr: ":" + config.AppConfig.Port,
		Handler: r  ,
	}
	
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func ()  {	
		err := server.ListenAndServe()

		if err != nil {
			log.Fatalln("Unable to run server", err)
		}
	}()

	<-done
	
	slog.Info("Initiating Server Shutdown Process")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Unable to shoutdown server", "ERROR", err)
	}


}