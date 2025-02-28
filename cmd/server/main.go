package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/tanjed/go-sso/apiservice"
	"github.com/tanjed/go-sso/internal/route"
)


func main() {
	apiServiceContainer := apiservice.NewApiService()
	apiServiceContainer.Boot()
	
	app := apiservice.GetApp()
	server := &http.Server{
		Addr: ":" + strconv.Itoa(app.Config.Port),
		Handler: route.NewRouter()  ,
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

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("Unable to shoutdown server", "ERROR", err)
	}
	
	apiServiceContainer.Destroy()
}

