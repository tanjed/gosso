package db

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/tanjed/go-sso/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DB_DRIVER_POSTGRES = "postgres"

type DB struct {
	Id int
	Conn *mongo.Client
}

var dbInstance *DB
var dbInstanceOnce sync.Once

func InitDB() *DB{

	dbInstanceOnce.Do(func ()  {
			fmt.Println(getConnectionString())
		client, err := mongo.Connect(options.Client().ApplyURI(getConnectionString()))

		if err != nil {
			slog.Error("Error creating mongo session:", "error", err)
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		if err:= client.Ping(ctx, nil); err != nil {
			slog.Error("Error creating mongo session:", "error", err)
			log.Fatal(err)
		}

		dbInstance = &DB{
			Id: rand.Intn(101),
			Conn: client,
		}
			
	})

	if dbInstance.Conn == nil {
		log.Fatal("Unable to establish DB connection") 
	}
	fmt.Println("DB CONNECTION ID:", dbInstance.Id)
	return dbInstance
}

func CloseDB() {
	if dbInstance != nil && dbInstance.Conn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()
		dbInstance.Conn.Disconnect(ctx)
	}
}

func getConnectionString() string {
	return "mongodb://"+config.AppConfig.DB_USER+":"+config.AppConfig.DB_PASSWORD+"@"+config.AppConfig.DB_HOST+":"+strconv.Itoa(config.AppConfig.DB_PORT)
}